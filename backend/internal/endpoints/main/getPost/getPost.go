package getPost

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/zpx64/supreme-octopus/internal/auth"
	"github.com/zpx64/supreme-octopus/internal/db"
	"github.com/zpx64/supreme-octopus/internal/model"
	"github.com/zpx64/supreme-octopus/internal/utils"
	"github.com/zpx64/supreme-octopus/internal/vars"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/ssleert/limiter"
)

var (
	// api endpoint like /put
	name   string
	logger zerolog.Logger
	limit  *limiter.Limiter[string]
	dbConn *pgxpool.Pool
)

// TODO: disable attachments in model.PostArticle
type Input struct {
	AccessToken string `json:"access_token" validate:"required,min=5,max=100"`
	Offset      uint   `json:"offset"`
	Limit       uint   `json:"limit"        validate:"required,max=35"`
}

type Post struct {
	Nickname             string           `json:"nickname"`
	AvatarImg            string           `json:"avatar_img"`
	Id                   int              `json:"id"`
	CreationDate         time.Time        `json:"creation_date"`
	Type                 model.Post       `json:"type"`
	Body                 string           `json:"body"`
	Attachments          []string         `json:"attachments"`
	VotesAmount          int              `json:"votes_amount"`
	VoteAction           model.VoteAction `json:"vote_action"`
	CommentsAmount       int              `json:"comments_amount"`
	IsCommentsDisallowed bool             `json:"is_comments_disallowed"`
}

type Output struct {
	Posts  []Post `json:"posts"`
	Err    string `json:"error"`
	Status int    `json:"-"`
}

func Start(n string, log *zerolog.Logger) error {
	var err error

	logger = *log
	name = n

	logger.Trace().Msg("creating req limiter")
	limit = limiter.New[string](vars.LimitPerHour, 3600, 2048, 4096, 20)

	logger.Trace().Msg("creating db connection")
	dbConn, err = pgxpool.New(context.Background(), db.GetConnString())
	if err != nil {
		logger.Error().
			Err(err).
			Msg("error with db connection")

		return err
	}

	logger.Trace().Msg("initing auth")
	err = auth.Init(context.Background(), logger)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("error with auth module initialization")

		return err
	}

	logger.Info().Msgf("%s endpoint started", name)
	return nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	//defer r.Body.Close()

	log := hlog.FromRequest(r)
	log.Info().Msg("connected")

	in := Input{}
	out := Output{
		Err:    "null",
		Status: http.StatusOK,
	}

	log.Trace().Interface("in", in).Send()

	defer func() {
		utils.WriteJsonAndStatusInRespone(w, &out, out.Status)
	}()

	var err error
	out.Status, err = utils.EndPointPrerequisites(
		log, w, r, limit, &in,
	)
	if err != nil {
		log.Warn().Err(err).Msg("preresquisites error")

		out.Err = err.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	accessTokenUint, err := strconv.ParseUint(in.AccessToken, 10, 64)
	if err != nil {
		log.Warn().Err(err).Msg("unsigned integer conversion error")

		out.Err = vars.ErrIncorrectUintValue.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	userId, err := auth.GetUserIdByAccessTokenWithDefaultToken(accessTokenUint)
	if err != nil {
		log.Warn().Err(err).Msg("error with access token")

		out.Err = err.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	ctx, cancel := context.WithTimeout(
		r.Context(),
		time.Duration(vars.TimeoutSeconds)*time.Second,
	)
	defer cancel()

	var (
		posts []model.UserNPost
	)
	err = dbConn.AcquireFunc(ctx,
		func(c *pgxpool.Conn) error {
			posts, err = db.ListPosts(ctx, c,
				in.Offset, in.Limit,
			)
			return err
		},
	)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("an error with database")

		out.Err = vars.ErrWithDb.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	voteActions := make([]model.VoteAction, len(posts))
	if userId != 0 {
		err = dbConn.AcquireFunc(ctx,
			func(c *pgxpool.Conn) error {
				// TODO: optimize without loop
				for i, e := range posts {
					voteAction, err := db.IsPostVoted(ctx, c,
						userId, e.Post.PostId,
					)
					if err != nil {
						return err
					}
					voteActions[i] = voteAction
				}
				return nil
			},
		)
		if err != nil {
			log.Warn().
				Err(err).
				Msg("an error with database")

			out.Err = vars.ErrWithDb.Error()
			out.Status = http.StatusInternalServerError
			return
		}
	}

	out.Posts = make([]Post, len(posts))
	for i, e := range posts {
		out.Posts[i] = Post{
			Nickname:             e.User.Nickname,
			AvatarImg:            e.User.AvatarImg,
			Id:                   e.Post.PostId,
			CreationDate:         e.Post.CreationDate,
			Type:                 e.Post.PostType,
			Body:                 e.Post.Body,
			Attachments:          e.Post.Attachments,
			VotesAmount:          e.Post.VotesAmount,
			VoteAction:           voteActions[i],
			CommentsAmount:       e.Post.CommentsAmount,
			IsCommentsDisallowed: e.Post.IsCommentsDisallowed,
		}
	}

	log.Debug().
		Interface("input_json", in).
		Interface("output_json", out).
		Send()
}

func Stop() error {
	if dbConn != nil {
		dbConn.Close()
	}
	logger.Info().Msgf("%s endpoint stoped", name)
	return nil
}

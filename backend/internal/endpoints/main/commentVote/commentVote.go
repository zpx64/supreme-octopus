package commentVote

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

type Input struct {
	AccessToken string           `json:"access_token"  validate:"required,max=256"`
	CommentId   int              `json:"comment_id"    validate:"required"`
	Action      model.VoteAction `json:"action"        validate:"required,min=1,max=2"`
}

type Output struct {
	LikeId int    `json:"like_id"`
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

	userId, err := auth.GetUserIdByAccessToken(accessTokenUint)
	if err != nil {
		log.Warn().Err(err).Msg("error with access token")

		out.Err = err.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	currentTime := time.Now()

	ctx, cancel := context.WithTimeout(
		r.Context(),
		time.Duration(vars.TimeoutSeconds)*time.Second,
	)
	defer cancel()

	var (
		id int
	)
	err = dbConn.AcquireFunc(ctx,
		func(c *pgxpool.Conn) error {
			id, err = db.VoteComment(ctx, c,
				userId,
				in.CommentId,
				in.Action,
				currentTime,
			)
			return err
		},
	)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("an error with database")

		if err == vars.ErrNotInDb {
			out.Err = err.Error()
		} else {
			out.Err = vars.ErrWithDb.Error()
		}
		out.Status = http.StatusInternalServerError
		return
	}

	out.LikeId = id

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

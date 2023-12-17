package reg

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/zpx64/supreme-octopus/internal/db"
	"github.com/zpx64/supreme-octopus/internal/model"
	"github.com/zpx64/supreme-octopus/internal/utils"
	"github.com/zpx64/supreme-octopus/internal/vars"

	"github.com/zpx64/supreme-octopus/pkg/cryptograph"
	"github.com/zpx64/supreme-octopus/pkg/valid"

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

// TODO: add validation for only ENGLISH alphabet in fields
type Input struct {
	Nickname string  `json:"nickname"           validate:"required,min=3,max=256"`
	Name     *string `json:"name,omitempty"`
	Surname  *string `json:"surname,omitempty"`
	Email    string  `json:"email"              validate:"required,min=5,max=256"`
	Password string  `json:"password"           validate:"required,min=6,max=256"`
}

type Output struct {
	WritedId int    `json:"writed_id"`
	Err      string `json:"error"`
	Status   int    `json:"-"`
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

	// TODO: rewrite with normal validation
	if in.Name != nil || in.Surname != nil {
		if in.Name == nil {
			log.Warn().Err(err).Msg("validation error")

			out.Err = errors.Join(
				vars.ErrOnValidation,
				errors.New("Name is not defined."),
			).Error()
			out.Status = http.StatusUnprocessableEntity
			return
		}

		if in.Surname == nil {
			log.Warn().Err(err).Msg("validation error")

			out.Err = errors.Join(
				vars.ErrOnValidation,
				errors.New("Surname is not defined."),
			).Error()
			out.Status = http.StatusUnprocessableEntity
			return
		}

		if len(*in.Name) < 3 || len(*in.Name) > 256 ||
			len(*in.Surname) < 3 || len(*in.Surname) > 256 {
			log.Warn().Err(err).Msg("validation error")

			out.Err = errors.Join(
				vars.ErrOnValidation,
				errors.New("Len of surname or name is out of range."),
			).Error()
			out.Status = http.StatusUnprocessableEntity
			return
		}
	}
	// *==========================================*

	if !valid.IsEmail(in.Email) {
		log.Warn().
			Str("email", in.Email).
			Msg("email is incorrect")

		out.Err = vars.ErrEmailIncorrect.Error()
		out.Status = http.StatusUnprocessableEntity
		return
	}

	localPow, err := cryptograph.GenRandPow(vars.PowLen)
	if err != nil {
		log.Warn().Err(err).Msg("error on random pow generation")

		out.Err = err.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	hashedPassword := utils.HashPassWithPows(in.Password, localPow)

	// TODO: rewrite with global contant
	dbModel := model.UserNCred{
		User: model.User{
			Nickname:     in.Nickname,
			AvatarImg:    "default",
			Name:         in.Name,
			Surname:      in.Surname,
			CreationDate: time.Now(),
		},
		Credentials: model.UserCredentials{
			Email:    in.Email,
			Password: hashedPassword,
			Pow:      localPow,
		},
	}

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
			id, err = db.CreateUser(ctx, c, &dbModel)
			return err
		},
	)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("an error with database")

		if err == vars.ErrAlreadyInDb {
			out.Err = vars.ErrAlreadyInDb.Error()
		} else {
			out.Err = vars.ErrWithDb.Error()
		}
		out.Status = http.StatusInternalServerError
		return
	}

	out.WritedId = id

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

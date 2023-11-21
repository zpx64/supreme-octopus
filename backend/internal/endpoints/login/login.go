package login

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

type Input struct {
	Email    string `json:"email"     validate:"required,min=3,max=256"`
	Password string `json:"password"  validate:"required,min=6,max=256"`
	DeviceId string `json:"device_id" validate:"required,min=6"`
}

type Output struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Err          string `json:"error"`
	Status       int    `json:"-"`
}

// TODO: incorrect error on empty json parsing
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

	err = auth.Init(context.Background())
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
	defer r.Body.Close()

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

	if !valid.IsEmail(in.Email) {
		log.Warn().
			Str("email", in.Email).
			Msg("email is incorrect")

		out.Err = vars.ErrEmailIncorrect.Error()
		out.Status = http.StatusUnprocessableEntity
		return
	}

	ctx, cancel := context.WithTimeout(
		r.Context(),
		time.Duration(vars.TimeoutSeconds)*time.Second,
	)
	defer cancel()

	var (
		credentials model.UserCredentials
	)
	err = dbConn.AcquireFunc(ctx,
		func(c *pgxpool.Conn) error {
			credentials, err = db.GetCredentialsByEmail(ctx, c, in.Email)
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

	passwordCorrect := (credentials.Password == utils.HashPassWithPows(
		in.Password, credentials.Pow,
	))
	if !passwordCorrect {
		log.Warn().
			Msg("user password is incorrect")

		out.Err = vars.ErrPassIncorrect.Error()
		out.Status = http.StatusForbidden
		return
	}

	deviceIdHash, err := auth.HashDeviceId(in.DeviceId)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("cant hash device id")

		out.Err = err.Error()
		out.Status = http.StatusUnprocessableEntity
		return
	}

	accessToken, refreshToken, err := auth.GenTokensPair(
		credentials.Id,
		deviceIdHash,
	)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("cant generate new tokens pair")

		out.Err = err.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	out.AccessToken = strconv.FormatUint(accessToken, 10)
	out.RefreshToken = refreshToken

	log.Debug().
		Interface("input_json", in).
		Interface("credentials", credentials).
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

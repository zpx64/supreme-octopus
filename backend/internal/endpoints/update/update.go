package update

import (
	"context"
	"net/http"

	"github.com/ssleert/tzproj/internal/db"
	"github.com/ssleert/tzproj/internal/utils"
	"github.com/ssleert/tzproj/internal/vars"

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

type input db.All
type output struct {
	Status int    `json:"-"`
	Err    string `json:"error"`
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
	defer r.Body.Close()

	log := hlog.FromRequest(r)
	log.Info().Msg("connected")

	in := input{}
	out := output{Err: "null"}

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
		return
	}

	log.Trace().
		Int("people_id", in.P.Id).
		Msg("updating data in db")

	err = dbConn.AcquireFunc(context.Background(), func(c *pgxpool.Conn) error {
		err = db.ReplaceAll(c, db.All(in))
		return err
	})
	if err == vars.ErrNotInDb {
		log.Warn().
			Err(err).
			Msg("data not in db")

		out.Err = vars.ErrNotInDb.Error()
		out.Status = http.StatusInternalServerError
		return
	}
	if err != nil {
		log.Warn().
			Err(err).
			Msg("an error with database")

		out.Err = vars.ErrWithDb.Error()
		out.Status = http.StatusInternalServerError
		return
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

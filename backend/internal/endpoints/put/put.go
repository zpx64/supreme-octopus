package put

import (
	"context"
	"net/http"

	"github.com/ssleert/tzproj/internal/conversions"
	"github.com/ssleert/tzproj/internal/db"
	"github.com/ssleert/tzproj/internal/utils"
	"github.com/ssleert/tzproj/internal/vars"

	"github.com/ssleert/tzproj/pkg/agify"
	"github.com/ssleert/tzproj/pkg/genderize"
	"github.com/ssleert/tzproj/pkg/nationalize"

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

	agifyClient       *agify.Client
	genderizeClient   *genderize.Client
	nationalizeClient *nationalize.Client
)

type input struct {
	Change     bool    `json:"change"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Patronymic *string `json:"patronymic,omitempty"`
}

type output struct {
	Status   int    `json:"-"`
	WritedId int    `json:"writed_id"`
	Err      string `json:"error"`
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

	logger.Trace().Msg("creating agify.io, genderize.io and nationalize.io clients")

	agifyClient, err = agify.New()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("agify.io client creation failed")

		return err
	}

	genderizeClient, err = genderize.New()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("genderize.io client creation failed")

		return err
	}

	nationalizeClient, err = nationalize.New()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("nationalize.io client creation failed")

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

	var (
		ageChan    = make(chan utils.Result[agify.Output])
		genderChan = make(chan utils.Result[genderize.Output])
		nationChan = make(chan utils.Result[nationalize.Output])
	)
	go func() {
		log.Trace().Msg("getting agify.io data")

		o, err := agifyClient.Get(in.Name)
		ageChan <- utils.Result[agify.Output]{o, err}

		close(ageChan)
	}()
	go func() {
		log.Trace().Msg("getting genderize.io data")

		o, err := genderizeClient.Get(in.Name)
		genderChan <- utils.Result[genderize.Output]{o, err}

		close(genderChan)
	}()
	go func() {
		log.Trace().Msg("getting nationalize.io data")

		o, err := nationalizeClient.Get(in.Name)
		nationChan <- utils.Result[nationalize.Output]{o, err}

		close(nationChan)
	}()

	ageResult := <-ageChan
	genderResult := <-genderChan
	nationResult := <-nationChan
	if ageResult.Err != nil {
		log.Warn().
			Err(err).
			Msg("agify.io error")

		out.Err = vars.ErrWithExternalApi.Error()
		out.Status = http.StatusInternalServerError
		return
	}
	if genderResult.Err != nil {
		log.Warn().
			Err(err).
			Msg("genderize.io error")

		out.Err = vars.ErrWithExternalApi.Error()
		out.Status = http.StatusInternalServerError
		return
	}
	if nationResult.Err != nil {
		log.Warn().
			Err(err).
			Msg("nationalize.io error")

		out.Err = vars.ErrWithExternalApi.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	age := ageResult.Val
	gender := genderResult.Val
	nation := nationResult.Val

	allData := db.All{
		P: db.People{
			Name:       in.Name,
			Surname:    in.Surname,
			Patronymic: in.Patronymic,
			Age:        age.Age,
		},
		G: db.Gender{
			Gender:      gender.Gender,
			Probability: gender.Probability,
		},
		N: conversions.CountriesToNationalization(
			nation.Countries,
		),
	}

	log.Trace().
		Interface("all_data", allData).
		Msg("writing data to db")

	id := 0
	err = dbConn.AcquireFunc(context.Background(), func(c *pgxpool.Conn) error {
		id, err = db.InsertAll(c, allData)
		return err
	})
	if err == vars.ErrAlreadyInDb {
		log.Warn().
			Err(err).
			Msg("data already in db")

		out.Err = err.Error()
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

	out.WritedId = id

	log.Debug().
		Interface("input_json", in).
		Interface("output_json", out).
		Bool("is_patronymic", in.Patronymic != nil).
		Send()
}

func Stop() error {
	if dbConn != nil {
		dbConn.Close()
	}
	logger.Info().Msgf("%s endpoint stoped", name)
	return nil
}

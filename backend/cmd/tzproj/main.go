package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ssleert/tzproj/internal/db"
	"github.com/ssleert/tzproj/internal/vars"

	// endpoints
	"github.com/ssleert/tzproj/internal/endpoints/del"
	"github.com/ssleert/tzproj/internal/endpoints/get"
	"github.com/ssleert/tzproj/internal/endpoints/put"
	"github.com/ssleert/tzproj/internal/endpoints/update"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

type endPoint struct {
	name    string
	start   func(string, *zerolog.Logger) error
	handler func(http.ResponseWriter, *http.Request)
	stop    func() error
}

var (
	endPoints = [...]endPoint{
		{"/put", put.Start, put.Handler, put.Stop},
		{"/del", del.Start, del.Handler, del.Stop},
		{"/update", update.Start, update.Handler, update.Stop},
		{"/get", get.Start, get.Handler, get.Stop},
	}
	logger zerolog.Logger
)

func main() {
	var (
		logFile     io.Writer
		loggerLevel zerolog.Level
		err         error
	)
	if vars.LogStdout {
		logFile = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
	} else {
		logFile, err = os.OpenFile(
			"./logs/"+time.Now().Format("2006_01_02-15:04:05")+".log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0666,
		)
		if err != nil {
			fmt.Println("log file creation failed")
			os.Exit(1)
		}
	}
	if vars.DebugMode {
		loggerLevel = zerolog.TraceLevel
	} else {
		loggerLevel = zerolog.InfoLevel
	}

	logger = zerolog.New(logFile).
		Level(loggerLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	logger.Info().Msg("started")
	vars.PrintVars(&logger)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		s := <-signalChannel
		switch s {
		case syscall.SIGHUP:
			logger.Error().Msg("SIGHUP received")
			stop()
			os.Exit(0)
		case syscall.SIGINT:
			logger.Error().Msg("SIGINT received")
			stop()
			os.Exit(0)
		case syscall.SIGTERM:
			logger.Error().Msg("SIGTERM received")
			stop()
			os.Exit(0)
		case syscall.SIGQUIT:
			logger.Error().Msg("SIGQUIT received")
			stop()
			os.Exit(0)
		default:
			logger.Error().Msg("unknown SIG received")
			stop()
			os.Exit(1)
		}
	}()

	migrationsDir, err := db.GetMigrationsDir()
	logger.Trace().
		Err(err).
		Str("migrations_dir", migrationsDir).
		Send()

	logger.Trace().Msg("making migrations")
	err = db.MakeMigrations(&logger)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("db migration failed")
	}

	// sorry for sooo big chain of append
	// i really want to make it clear(
	mux := http.NewServeMux()
	h := alice.New().
		Append(hlog.NewHandler(logger)).
		Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Send()
		})).
		Append(hlog.RemoteAddrHandler("ip")).
		Append(hlog.UserAgentHandler("user_agent")).
		Append(hlog.RefererHandler("referer")).
		Append(hlog.RequestIDHandler("req_id", "Request-Id")).
		Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			mux.ServeHTTP(w, r)
		}))

	for _, e := range endPoints {
		mux.HandleFunc(vars.EndPointPrefix+e.name, e.handler)
	}

	start()
	defer stop()
	logger.Info().Msgf("http server starting on %s port", vars.HttpPort)
	err = http.ListenAndServe(":"+vars.HttpPort, h)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("http server start failed")
	}
	logger.Info().Msg("exiting")
}

func start() {
	for i, e := range endPoints {
		logger.Info().Msgf("%s start called", e.name)
		err := e.start(e.name, &logger)
		if err != nil {
			for j := 0; j <= i; j++ {
				logger.Info().Msgf("%s stop called", endPoints[j].name)
				if err := endPoints[j].stop(); err != nil {
					logger.Error().Err(err)
				}
			}
			logger.Fatal().Err(err)
		}
	}
}

func stop() {
	for _, e := range endPoints {
		logger.Info().Msgf("%s stop called", e.name)
		err := e.stop()
		if err != nil {
			logger.Fatal().Err(err)
		}
	}
}

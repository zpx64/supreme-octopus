package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zpx64/supreme-octopus/internal/db"
	"github.com/zpx64/supreme-octopus/internal/vars"

	// endpoints
	"github.com/zpx64/supreme-octopus/internal/endpoints/main/commentNew"
	"github.com/zpx64/supreme-octopus/internal/endpoints/main/commentVote"
	"github.com/zpx64/supreme-octopus/internal/endpoints/main/delVote"
	"github.com/zpx64/supreme-octopus/internal/endpoints/main/delCommentVote"
	"github.com/zpx64/supreme-octopus/internal/endpoints/main/getPost"
	"github.com/zpx64/supreme-octopus/internal/endpoints/main/listComments"
	"github.com/zpx64/supreme-octopus/internal/endpoints/main/listPosts"
	"github.com/zpx64/supreme-octopus/internal/endpoints/main/login"
	"github.com/zpx64/supreme-octopus/internal/endpoints/main/postImage"
	"github.com/zpx64/supreme-octopus/internal/endpoints/main/postNew"
	"github.com/zpx64/supreme-octopus/internal/endpoints/main/postVote"
	"github.com/zpx64/supreme-octopus/internal/endpoints/main/refresh"
	"github.com/zpx64/supreme-octopus/internal/endpoints/main/reg"
	"github.com/zpx64/supreme-octopus/internal/endpoints/main/test"

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
		{"/test", test.Start, test.Handler, test.Stop},
		{"/reg", reg.Start, reg.Handler, reg.Stop},
		{"/login", login.Start, login.Handler, login.Stop},
		{"/refresh", refresh.Start, refresh.Handler, refresh.Stop},
		{"/post_new", postNew.Start, postNew.Handler, postNew.Stop},
		{"/comment_new", commentNew.Start, commentNew.Handler, commentNew.Stop},
		{"/post_image", postImage.Start, postImage.Handler, postImage.Stop},
		{"/list_posts", listPosts.Start, listPosts.Handler, listPosts.Stop},
		{"/get_post", getPost.Start, getPost.Handler, getPost.Stop},
		{"/list_comments", listComments.Start, listComments.Handler, listComments.Stop},
		{"/post_vote", postVote.Start, postVote.Handler, postVote.Stop},
		{"/comment_vote", commentVote.Start, commentVote.Handler, commentVote.Stop},
		{"/del_vote", delVote.Start, delVote.Handler, delVote.Stop},
		{"/del_comment_vote", delCommentVote.Start, delCommentVote.Handler, delCommentVote.Stop},
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
		// TODO: test json logging to file
		//       maybe port to something like logstash
		//       but more lightweight and simple
		logFile, err = os.OpenFile(
			vars.LogPath+"/backend_main"+time.Now().Format("2006_01_02-15:04:05")+".log",
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

	server := &http.Server{
		Addr:           ":" + vars.HttpPort,
		Handler:        h,
		ReadTimeout:    time.Duration(vars.ReadTimeoutSeconds) * time.Second,
		WriteTimeout:   time.Duration(vars.WriteTimeoutSeconds) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger.Info().Msgf("http server starting on %s port", vars.HttpPort)
	err = server.ListenAndServe()
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
			logger.Error().Err(err).Msg("an error on module init")
			for j := 0; j < i; j++ {
				logger.Info().Msgf("%s stop called", endPoints[j].name)
				if err := endPoints[j].stop(); err != nil {
					logger.Error().Err(err)
				}
			}
			logger.Fatal().Err(err)
			os.Exit(1)
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

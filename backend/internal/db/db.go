package db

import (
	"fmt"
	"os"

	"github.com/ssleert/tzproj/internal/vars"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
)

type zerologLogger struct {
	l *zerolog.Logger
}

func (log zerologLogger) Printf(f string, v ...any) {
	log.l.Debug().Msgf("MIGRATION: "+f[:len(f)-1], v...)
}

func (log zerologLogger) Verbose() bool {
	return true
}

func GetConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:5432/postgres?%s",
		vars.PostgresUser,
		vars.PostgresPassword,
		vars.PostgresDbUrl,
		vars.PostgresConnFlags,
	)
}

func GetMigrationsDir() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path + "/" + vars.MigrationsPath, nil
}

func MakeMigrations(log *zerolog.Logger) error {
	path, err := GetMigrationsDir()
	if err != nil {
		return err
	}
	m, err := migrate.New(
		"file://"+path,
		GetConnString(),
	)
	if err != nil {
		return err
	}
	m.Log = zerologLogger{log}
	m.Log.Printf("created migration client ")

	if vars.PostgresForceDrop {
		m.Log.Printf("force dropping db ")
		err = m.Down()
		if err != nil && err != migrate.ErrNoChange {
			return err
		}
	}
	m.Log.Printf("uping migrations ")
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

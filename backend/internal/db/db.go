package db

import (
	"fmt"
	"os"

	"github.com/zpx64/supreme-octopus/internal/db/migrations"
	"github.com/zpx64/supreme-octopus/internal/vars"

	"github.com/rs/zerolog"
)

type zerologLogger struct {
	l *zerolog.Logger
}

func (log zerologLogger) Printf(f string, v ...any) {
	log.l.Debug().Msgf("MIGRATION: "+f, v...)
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

// TODO: rewrite without idiot logging interface
//
//	it written to integrate zerolog with go-migrate
//	but now we use tern and it uneeded
func MakeMigrations(log *zerolog.Logger) error {
	logger := zerologLogger{log}

	migrator, err := migrations.NewMigrator(GetConnString())
	if err != nil {
		return err
	}
	logger.Printf("created migration client")

	if vars.PostgresForceDrop {
		logger.Printf("force dropping all migrations")
		err = migrator.MigrateTo(0)
		if err != nil {
			return err
		}
	}

	now, exp, info, err := migrator.Info()
	if err != nil {
		return err
	}
	logger.Printf("getted migration info")

	logger.Printf("checking migration state")
	if now < exp {
		logger.Printf("current state: %s", info)

		err = migrator.Migrate()
		if err != nil {
			return err
		}
		logger.Printf("migration successful")
		return nil
	}

	logger.Printf("migration not needed")
	return nil
}

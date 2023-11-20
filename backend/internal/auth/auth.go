package auth

import (
	"context"

	"github.com/zpx64/supreme-octopus/internal/db"
	"github.com/zpx64/supreme-octopus/internal/vars"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ssleert/mumap"
)

type token struct {
	date     int64
	id       int
	deviceId uint64
}

type tokenMaps struct {
	accessTokens  mumap.Map[uint64, token]
	refreshTokens mumap.Map[string, token]
}

var (
	dbConnPool *pgxpool.Pool

	tokens = tokenMaps{
		accessTokens:  mumap.New[uint64, token](vars.DefaultMapSize),
		refreshTokens: mumap.New[string, token](vars.DefaultMapSize),
	}
)

// i dont really understand why we need context here
func Init(ctx context.Context) error {
	var err error

	dbConnPool, err = pgxpool.New(
		ctx, db.GetConnString(),
	)
	if err != nil {
		return err
	}

	return nil
}

package auth

import (
	"context"

	"github.com/zpx64/supreme-octopus/internal/db"
	"github.com/zpx64/supreme-octopus/internal/vars"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ssleert/mumap"
	"github.com/rs/zerolog"
)

type refreshToken struct {
	dbId int
	date int64
}

type tokenMaps struct {
	accessTokens  mumap.Map[uint64, int64]
	refreshTokens mumap.Map[string, refreshToken]
}

var (
	inited bool
	logger zerolog.Logger
	dbConnPool *pgxpool.Pool

	tokens = tokenMaps{
		accessTokens:  mumap.New[uint64, int64](vars.DefaultMapSize),
		refreshTokens: mumap.New[string, refreshToken](vars.DefaultMapSize),
	}
)

// i dont really understand why we need context here
func Init(ctx context.Context, log zerolog.Logger) error {
	if inited {
		return nil
	}

	logger = log
	
	var err error
	dbConnPool, err = pgxpool.New(
		ctx, db.GetConnString(),
	)
	if err != nil {
		return err
	}

	inited = true
	return nil
}

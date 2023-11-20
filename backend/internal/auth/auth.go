package auth

import (
	"errors"
	"context"

	"github.com/zpx64/supreme-octopus/internal/db"
	"github.com/zpx64/supreme-octopus/internal/vars"

	"github.com/ssleert/mumap"
	"github.com/jackc/pgx/v5/pgxpool"
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
	ErrAccessTNotExpired = errors.New("Access token is not expired.")
	ErrRefreshTExpired   = errors.New("Refresh token is expired.")
	ErrRefreshTNotFound  = errors.New("Refresh token doesnt found.")
	ErrAccessTNotFound   = errors.New("Access token doesnt found.")
	ErrAccessTExpired    = errors.New("Access token expired.")
	ErrDeviceIdLenIsBiggerThanExpected = errors.New("Device id len is bigger than expected.")

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

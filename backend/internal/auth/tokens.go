package auth

import (
	"context"
	"time"
	"strconv"

	"github.com/zpx64/supreme-octopus/internal/vars"
	"github.com/zpx64/supreme-octopus/internal/db"
	"github.com/zpx64/supreme-octopus/internal/model"
	"github.com/zpx64/supreme-octopus/pkg/cryptograph"

	"github.com/cespare/xxhash"
	"github.com/jackc/pgx/v5/pgxpool"
)

// gen new hash by device id string
func HashDeviceId(deviceId string) (uint64, error) {
	if len(deviceId) > vars.MaxDeviceIdLen {
		return 0, vars.ErrAuthDeviceIdLenIsBiggerThanExpected
	}
	// i really dont trust what frontenders send for me)
	return xxhash.Sum64String(deviceId), nil
}

// check is access token correct and not expired
func ValidateAccessToken(tkn uint64) error {
	token, ok := tokens.accessTokens.Get(tkn)
	if !ok {
		return vars.ErrAuthAccessTNotFound
	}

	timeNow := time.Now().Unix()
	if timeNow-token.date >= vars.AccessTokenLifeSeconds {
		return vars.ErrAuthAccessTExpired
	}
	return nil
}

func GetUserIdByAccessToken(tkn uint64) (int, error) {
	err := ValidateAccessToken(tkn)
	if err != nil {
		return 0, err
	}
	
	token, _ := tokens.accessTokens.Get(tkn)

	return token.userId, nil
}

func GetTokens() (uint64, string, error) {	
	hash, err := cryptograph.GenRandHash() // access token
	if err != nil {
		return 0, "", err
	}
	uid, err := cryptograph.GenRandUuid() // refresh token (realization may be changed)
	if err != nil {
		return 0, "", err
	}

	return hash, uid, nil
}

// gen new tokens pair by device and user id's
func GenTokensPair(
	ctx context.Context,
	id int,
	deviceId uint64,
	userAgent string,
) (uint64, string, error) {

	hash, uid, err := GetTokens()
	if err != nil {
		return 0, "", err
	}

	timeNow := time.Now().Unix()

	tokenId := 0
	err = dbConnPool.AcquireFunc(ctx,
		func(c *pgxpool.Conn) error {
			var err error
			tokenId, err = db.InsertNewToken(
				ctx, c, model.UserToken{
					UserId: id,
					DeviceId: strconv.FormatUint(deviceId, 10),
					RefreshToken: uid,
					UserAgent: userAgent,
					TokenDate: timeNow,
				},
			)
			return err
		},
	)
	if err != nil {
		logger.Warn().Err(err).Msg("auth: error with db")
		return 0, "", vars.ErrWithDb
	}

	tokens.accessTokens.Set(hash, 
		accessToken{
			userId: id,
			date: timeNow,
		},
	)
	tokens.refreshTokens.Set(uid, 
		refreshToken{
			dbId: tokenId,
			date: timeNow,
		},
	)

	return hash, uid, nil
}

// func that generates new tokens pair
// and remove expired and useless tokens
// from db and hashmap cache
// TODO: rewrite with normal hashmap data update
func RegenTokensPair(
	ctx context.Context,
	access uint64,
	refresh string,
	userAgent string,
) (uint64, string, error) {
	timeNow := time.Now().Unix()

	token, ok := tokens.accessTokens.Get(access)
	if ok {
		if timeNow-token.date < vars.AccessTokenLifeSeconds {
			//println("skipped by your mommy)")
			return 0, "", vars.ErrAuthAccessTNotExpired
		}
		tokens.accessTokens.Del(access)
	}

	t, ok := tokens.refreshTokens.Get(refresh)
	if !ok {
		return 0, "", vars.ErrAuthRefreshTNotFound
	}
	if timeNow-t.date > vars.RefreshTokenLifeSeconds {
		return 0, "", vars.ErrAuthRefreshTExpired
	}

	tokens.refreshTokens.Del(refresh)

	hash, uid, err := GetTokens()
	if err != nil {
		return 0, "", err
	}

	timeNow = time.Now().Unix()
	err = dbConnPool.AcquireFunc(ctx,
		func(c *pgxpool.Conn) error {
			err := db.UpdateToken(
				ctx, c, model.UserToken{
					TokenId: t.dbId,
					RefreshToken: uid,
					TokenDate: timeNow,
				},
			)
			return err
		},
	)
	if err != nil {
		if err == vars.ErrNotInDb {
			return 0, "", err
		}

		logger.Warn().Err(err).Msg("auth: error with db")
		return 0, "", vars.ErrWithDb
	}

	tokens.accessTokens.Set(hash, 
		accessToken{
			userId: token.userId,
			date: timeNow,
		},	
	)
	tokens.refreshTokens.Set(uid, 
		refreshToken{
			dbId: t.dbId,
			date: timeNow,
		},
	)

	return hash, uid, nil
}

package auth

import (
	"time"

	"github.com/zpx64/supreme-octopus/internal/vars"
	"github.com/zpx64/supreme-octopus/pkg/cryptograph"

	"github.com/cespare/xxhash"
)

// gen new hash by device id string
func HashDeviceId(deviceId string) (uint64, error) {
	if len(deviceId) > vars.MaxDeviceIdLen {
		return 0, ErrDeviceIdLenIsBiggerThanExpected
	}
	// i really dont trust what frontenders send for me)
	return xxhash.Sum64String(deviceId), nil
}

// check is access token correct and not expired
func ValidateAccessToken(tkn uint64) (bool, error) {
	timeNow := time.Now().Unix()

	t, ok := tokens.accessTokens.Get(tkn)
	if !ok {
		return false, ErrAccessTNotFound
	}
	if timeNow-t.date >= vars.AccessTokenLifeSeconds {
		return false, ErrAccessTExpired
	}
	return true, nil
}

// gen new tokens pair by device and user id's
func GenTokensPair(
	id int,
	deviceId uint64,
) (uint64, string, error) {
	timeNow := time.Now().Unix()

	hash, err := cryptograph.GenRandHash() // access token
	if err != nil {
		return 0, "", err
	}
	uid, err := cryptograph.GenRandUuid() // refresh token (realization may be changed)
	if err != nil {
		return 0, "", err
	}

	// TODO: insert refresh token in db

	tokens.accessTokens.Set(
		hash,
		token{
			date:     timeNow,
			id:       id,
			deviceId: deviceId,
		},
	)
	tokens.refreshTokens.Set(
		uid,
		token{
			date:     timeNow,
			id:       id,
			deviceId: deviceId,
		},
	)

	return hash, uid, nil
}

// func that generates new tokens pair
// and remove expired and useless tokens
// from db and hashmap cache
func RegenTokensPair(
	access uint64,
	refresh string,
) (uint64, string, error) {
	timeNow := time.Now().Unix()

	t, ok := tokens.accessTokens.Get(access)
	if ok {
		if timeNow-t.date < vars.AccessTokenLifeSeconds {
			return 0, "", ErrAccessTNotExpired
		}
		tokens.accessTokens.Del(access)
	}

	t, ok = tokens.refreshTokens.Get(refresh)
	if !ok {
		return 0, "", ErrRefreshTNotFound
	}
	if timeNow-t.date > vars.RefreshTokenLifeSeconds {
		return 0, "", ErrRefreshTExpired
	}

	tokens.refreshTokens.Del(refresh)

	return GenTokensPair(t.id, t.deviceId)
}

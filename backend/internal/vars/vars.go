package vars

import (
	"errors"
	"os"
	"strconv"

	"github.com/caitlinelfring/go-env-default"
	"github.com/rs/zerolog"
)

var (
	ErrWithDb               = errors.New("1: An Error with db.")
	ErrBodyLenIsTooBig      = errors.New("2: Request body too len for processing.")
	ErrInputJsonIsIncorrect = errors.New("3: Input json is incorrect.")
	ErrBodyReadingFailed    = errors.New("4: Error while req body reading.")
	ErrFieldTooBig          = errors.New("5: Json field is too big.")
	ErrActionLimited        = errors.New("6: Action limited. Try another time.")
	ErrWithPowGen           = errors.New("7: Pow generation failed. Try another time.")
	ErrEmailIncorrect       = errors.New("8: Email pattern is incorrect. (not x@y.zz)")
	ErrPassIncorrect        = errors.New("9: Password is incorrect.")
	ErrAlreadyInDb          = errors.New("10: Already in db.")
	ErrNotInDb              = errors.New("11: Not in db.")
	ErrWithExternalApi      = errors.New("12: External api is inaccessible.")
	ErrIncorrectFilterKey   = errors.New("13: Invalid filter key.")
	ErrOnValidation         = errors.New("14: Validation error: ") // need to be wrapped
	ErrEmailNotFound        = errors.New("15: Email not found.")

	// values from env
	HttpPort            = env.GetDefault("HTTP_PORT", "9876")
	TimeoutSeconds      = env.GetIntDefault("TIMEOUT_SECONDS", 30)
	ReadTimeoutSeconds  = env.GetIntDefault("READ_TIMEOUT_SECONDS", 60)
	WriteTimeoutSeconds = env.GetIntDefault("WRITE_TIMEOUT_SECONDS", 80)
	EndPointPrefix      = env.GetDefault("API_PREFIX", "")
	LogStdout           = env.GetBoolDefault("LOG_STDOUT", true)
	DebugMode           = env.GetBoolDefault("DEBUG_MODE", true)
	PostgresUser        = env.GetDefault("POSTGRES_USER", "admin")
	PostgresPassword    = env.GetDefault("POSTGRES_PASSWORD", "admin")
	PostgresDbUrl       = env.GetDefault("POSTGRES_DB_URL", "postgres")
	PostgresConnFlags   = env.GetDefault("POSTGRES_CONN_FLAGS", "")
	PostgresForceDrop   = env.GetBoolDefault("POSTGRES_FORCE_DROP", false) // drop db before start
	GlobalPow           = env.GetDefault("GLOBAL_POW", "btwsarseniyshouldsuckmydicknballs")
	PowLen              = env.GetIntDefault("POW_LEN", 32)
	PowRightCat         = env.GetBoolDefault("POW_RIGHT_CAT", true)
	MaxBodyLen          = env.GetInt64Default("MAX_BODY_LEN", 16384)
	MaxDeviceIdLen      = env.GetIntDefault("MAX_DEVICE_ID_LEN", 256)
	LimitPerHour        = env.GetIntDefault("LIMIT_PER_HOUR", 60)
	DefaultMapSize      = env.GetIntDefault("DEFAULT_MAP_SIZE", 2048)
	MigrationsPath      = env.GetDefault("MIGRATIONS_PATH", "./db/migrations")
	EmailMaxLen         = env.GetIntDefault("EMAIL_MAX_LEN", 256)
	NicknameMaxLen      = env.GetIntDefault("NICKNAME_MAX_LEN", 256)
	PasswordMaxLen      = env.GetIntDefault("PASSWORD_MAX_LEN", 256)
	AccessTokenLifeSeconds  = env.GetInt64Default("ACCESS_TOKEN_LIFE_SECONDS", 3600)
	RefreshTokenLifeSeconds = env.GetInt64Default("REFRESH_TOKEN_LIFE_SECONDS", 1814400)
)

func PrintVars(log *zerolog.Logger) {
	log.Trace().
		Str("HttpPort", HttpPort).
		Int("TimeoutSeconds", TimeoutSeconds).
		Int("ReadTimeoutSeconds", ReadTimeoutSeconds).
		Int("WriteTimeoutSeconds", WriteTimeoutSeconds).
		Str("EndPointPrefix", EndPointPrefix).
		Bool("LogStdout", LogStdout).
		Bool("DebugMode", DebugMode).
		Str("PostgresUser", PostgresUser).
		Str("PostgresPassword", PostgresPassword).
		Str("PostgresDbUrl", PostgresDbUrl).
		Str("PostgresConnFlags", PostgresConnFlags).
		Bool("PostgresForceDrop", PostgresForceDrop).
		Int64("MaxBodyLen", MaxBodyLen).
		Int("LimitPerHour", LimitPerHour).
		Str("MigrationsPath", MigrationsPath).
		Send()
}

type convertable interface {
	int | int64 | float64 | string
}

// WHY FUCKING GOLANG CANT DO TYPE SWITCH
// ON TYPE PARAMETER WHY U CANT FUCKING ASSHOLE
// INSTANTIATE FUCKING COMPILE-TIME SWITCH
// WHAT THE HECK
//
// okey, now i am chill and
// any(def).(type) can work properly)
func CheckEnv[T convertable](s string, def T) T {
	b := os.Getenv(s)
	if b == "" {
		return def
	}

	var (
		result any
		err    error
	)
	switch any(def).(type) {
	case int:
		result, err = strconv.Atoi(b)
		if err != nil {
			return def
		}
	case int64:
		l, err := strconv.Atoi(b)
		if err != nil {
			return def
		}
		result = int64(l)
	case float64:
		result, err = strconv.ParseFloat(b, 64)
		if err != nil {
			return def
		}
	case string:
		result = b
	}
	return result.(T)
}

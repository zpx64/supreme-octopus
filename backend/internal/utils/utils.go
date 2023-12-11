package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/zpx64/supreme-octopus/internal/vars"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/ssleert/limiter"
)

// TODO: split functions in separate files

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
)

type Result[T any] struct {
	Val T
	Err error
}

// simply serialize to json and write struct as http response with status code
func WriteJsonAndStatusInRespone[T any](w http.ResponseWriter, j *T, status int) {
	w.WriteHeader(status)
	jsn, _ := json.Marshal(*j)
	w.Write(jsn)
}

// simply serialize to json and write struct as http response with status code
func WriteStringAndStatusInRespone(w http.ResponseWriter, j *string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(*j))
}

func GetAddrFromStr(addrNPort *string) string {
	return strings.Split(
		*addrNPort, ":",
	)[0]
}

func ValidateStruct[T any](
	log *zerolog.Logger,
	str *T,
) error {
	log.Trace().Msg("validating data in struct")
	err := validate.Struct(str)
	if err != nil {
		return errors.Join(vars.ErrOnValidation, err)
	}
	return nil
}

func LimitUserByRemoteAddr(
	log *zerolog.Logger,
	r *http.Request,
	limit *limiter.Limiter[string],
) error {
	log.Trace().Msg("checking req limiter")

	if !limit.Try(GetAddrFromStr(&r.RemoteAddr)) {
		log.Warn().Msg("action limited")

		return vars.ErrActionLimited
	}
	return nil
}

func CheckBodyLen(log *zerolog.Logger, r *http.Request) error {
	log.Trace().Msg("checking content len")

	if r.ContentLength > vars.MaxBodyLen {
		log.Warn().
			Int64("content_length", r.ContentLength).
			Int64("max_content_length", vars.MaxBodyLen).
			Msg("content length is too big")

		return vars.ErrBodyLenIsTooBig
	}
	return nil
}

func ReadAllBody(log *zerolog.Logger, r *http.Request, body *[]byte) error {
	var err error
	log.Trace().Msg("reading body")

	_, err = io.ReadFull(r.Body, *body)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("cant read all body")

		return vars.ErrBodyReadingFailed
	}
	return nil
}

func UnmarshalJson[T any](
	log *zerolog.Logger,
	body io.Reader,
	str *T,
) error {
	log.Trace().Msg("unmarshaling json")

	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	err := dec.Decode(str)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("cant unmarshal body to json")

		return vars.ErrInputJsonIsIncorrect
	}
	return nil
}

func BinaryEndPointPrerequisites(
	log *zerolog.Logger,
	r *http.Request,
	limit *limiter.Limiter[string],
	body *[]byte,
) (int, error) {
	err := LimitUserByRemoteAddr(log, r, limit)
	if err != nil {
		return http.StatusTooManyRequests, err
	}

	err = ReadAllBody(log, r, body)
	if err != nil {
		return http.StatusInsufficientStorage, err
	}

	return http.StatusOK, nil
}

func EndPointPrerequisitesWithoutMaxBodyLen[T any](
	log *zerolog.Logger,
	w http.ResponseWriter,
	r *http.Request,
	limit *limiter.Limiter[string],
	in *T,
) (int, error) {
	err := LimitUserByRemoteAddr(log, r, limit)
	if err != nil {
		return http.StatusTooManyRequests, err
	}

	err = UnmarshalJson(log, r.Body, in)
	if err != nil {
		return http.StatusUnprocessableEntity, err
	}

	err = ValidateStruct(log, in)
	if err != nil {
		return http.StatusUnprocessableEntity, err
	}

	return http.StatusOK, nil
}

func EndPointPrerequisites[T any](
	log *zerolog.Logger,
	w http.ResponseWriter,
	r *http.Request,
	limit *limiter.Limiter[string],
	in *T,
) (int, error) {
	err := LimitUserByRemoteAddr(log, r, limit)
	if err != nil {
		return http.StatusTooManyRequests, err
	}

	err = CheckBodyLen(log, r)
	if err != nil {
		return http.StatusRequestEntityTooLarge, err
	}

	err = UnmarshalJson(log, r.Body, in)
	if err != nil {
		return http.StatusUnprocessableEntity, err
	}

	err = ValidateStruct(log, in)
	if err != nil {
		return http.StatusUnprocessableEntity, err
	}

	return http.StatusOK, nil
}

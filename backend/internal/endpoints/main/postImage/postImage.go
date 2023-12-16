package postImage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zpx64/supreme-octopus/internal/auth"
	"github.com/zpx64/supreme-octopus/internal/imagesStore"
	"github.com/zpx64/supreme-octopus/internal/utils"
	"github.com/zpx64/supreme-octopus/internal/vars"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/ssleert/limiter"
)

var (
	// api endpoint like /put
	name   string
	logger zerolog.Logger

	limit *limiter.Limiter[string]
)

type Image struct {
	ContentType  string `json:"content_type"`
	EncodedImage string `json:"encoded_image"` // encoded in z85
}

type Input struct {
	AccessToken string  `json:"access_token" validate:"required,min=5,max=100"`
	Images      []Image `json:"images"`
}

type InputForImagesStore struct {
	Images []Image `json:"images"`
}

type Output struct {
	WritedIds []string `json:"writed_ids"`
	Err       string   `json:"error"`
	Status    int      `json:"-"`
}

func Start(n string, log *zerolog.Logger) error {
	logger = *log
	name = n

	logger.Trace().Msg("creating req limiter")
	limit = limiter.New[string](vars.LimitPerHour, 3600, 2048, 4096, 20)

	logger.Trace().Msg("initing auth")
	err := auth.Init(context.Background(), logger)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("error with auth module initialization")

		return err
	}

	logger.Info().Msgf("%s endpoint started", name)
	return nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	//defer r.Body.Close()

	log := hlog.FromRequest(r)
	log.Info().Msg("connected")

	in := Input{}
	out := Output{
		WritedIds: make([]string, 0, 32),
		Err:       "null",
		Status:    http.StatusOK,
	}

	defer func() {
		utils.WriteJsonAndStatusInRespone(w, &out, out.Status)
	}()

	var err error
	out.Status, err =
		utils.EndPointPrerequisitesWithoutMaxBodyLen(
			log, w, r, limit, &in,
		)
	if err != nil {
		log.Warn().Err(err).Msg("preresquisites error")

		out.Err = err.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	accessTokenUint, err := strconv.ParseUint(in.AccessToken, 10, 64)
	if err != nil {
		log.Warn().Err(err).Msg("unsigned integer conversion error")

		out.Err = vars.ErrIncorrectUintValue.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	err = auth.ValidateAccessToken(accessTokenUint)
	if err != nil {
		log.Warn().Err(err).Msg("error with access token")

		out.Err = err.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	microserviceInput := InputForImagesStore{
		Images: in.Images,
	}

	jsonBody, err := json.Marshal(microserviceInput)
	if err != nil {
		log.Warn().Err(err).Msg("error with json marshaling")

		out.Err = vars.ErrInternalJsonParsing.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	// TODO: add timeout context support
	resp, err := http.Post(
		fmt.Sprintf("http://%s/api/post_image", vars.ImagesStoreUrl), "text/json",
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		log.Warn().Err(err).Msg("an error with http post request")

		out.Err = vars.ErrInternalMicroserviceRequest.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	// TODO: rewrite without json unmarshal
	//       on plain write to response
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&out)
	if err != nil {
		log.Warn().Err(err).Msg("an error with json unmarshaling")

		out.Err = vars.ErrInternalJsonParsing.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	log.Debug().
		Interface("output_json", out).
		Send()
}

func Stop() error {
	err := imagesStore.Deinit()
	if err != nil {
		logger.Warn().Err(err).Msg("an error with images store")
	}
	logger.Info().Msgf("%s endpoint stoped", name)
	return nil
}

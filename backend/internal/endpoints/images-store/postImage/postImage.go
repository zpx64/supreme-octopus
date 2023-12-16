package postImage

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zpx64/supreme-octopus/internal/imagesStore"
	"github.com/zpx64/supreme-octopus/internal/utils"
	"github.com/zpx64/supreme-octopus/internal/vars"

	"github.com/nofeaturesonlybugs/z85"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var (
	// api endpoint like /put
	name   string
	logger zerolog.Logger
)

type Image struct {
	ContentType  string `json:"content_type"  validate:"required,min=5"`
	EncodedImage string `json:"encoded_image" validate:"required,min=10"` // encoded in z85
}

type Input struct {
	Images []Image `json:"images" validate:"required,min=1,max=8,dive"`
}

type Output struct {
	WritedIds []string `json:"writed_ids"`
	Err       string   `json:"error"`
	Status    int      `json:"-"`
}

func Start(n string, log *zerolog.Logger) error {
	logger = *log
	name = n

	logger.Info().Msg("initing images store")
	err := imagesStore.Init(log)
	if err != nil {
		logger.Warn().Err(err).Msg("an error with images store")
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
		utils.EndPointPrerequisitesWithoutLimiterAndMaxBodyLen(
			log, w, r, &in,
		)
	if err != nil {
		log.Warn().Err(err).Msg("preresquisites error")

		out.Err = err.Error()
		out.Status = http.StatusInternalServerError
		return
	}

	for i, image := range in.Images {
		log.Trace().Int("input_size", len(in.Images[i].EncodedImage)).Msg("input size of images")

		decodedImage, err := z85.PaddedDecode(image.EncodedImage)
		if err != nil {
			log.Warn().Err(err).Msg("z85 decoding error")

			out.Err = errors.Join(
				vars.ErrZ85Incorrect,
				fmt.Errorf("on %d image", i),
			).Error()
			out.Status = http.StatusUnprocessableEntity
			return
		}

		writedImageHash, err := imagesStore.PostImageToStore(decodedImage, image.ContentType)
		if err != nil {
			log.Warn().Err(err).Msg("image posting to store error")

			out.Err = errors.Join(
				vars.ErrImageUpload,
				fmt.Errorf("on %d image", i),
			).Error()
			out.Status = http.StatusUnprocessableEntity
			return
		}
		out.WritedIds = append(out.WritedIds, strconv.FormatUint(writedImageHash, 10))
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

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
	WritedIds   []string `json:"writed_ids"`
	ImageErrors []string `json:"image_errors"`
	Err         string   `json:"error"`
	Status      int      `json:"-"`
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
		Err:    "null",
		Status: http.StatusOK,
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

	results := make([]chan utils.Result[uint64], len(in.Images))
	for i := range results {
		results[i] = make(chan utils.Result[uint64])
	}

	for i, image := range in.Images {
		log.Trace().Int("input_size", len(in.Images[i].EncodedImage)).Msg("input size of images")

		go func(i int, image Image, res chan utils.Result[uint64]) {
			log.Trace().Msg("decoding z85")
			decodedImage, err := z85.PaddedDecode(image.EncodedImage)
			if err != nil {
				log.Warn().Err(err).Msg("z85 decoding error")

				res <- utils.Result[uint64]{
					Err: errors.Join(
						vars.ErrZ85Incorrect,
						fmt.Errorf("on %d image", i),
					),
				}
				return
			}

			log.Trace().Msg("posting image to store")
			writedImageHash, err := imagesStore.PostImageToStore(decodedImage, image.ContentType)
			if err != nil {
				log.Warn().Err(err).Msg("image posting to store error")

				res <- utils.Result[uint64]{
					Err: errors.Join(
						vars.ErrImageUpload,
						fmt.Errorf("on %d image", i),
					),
				}
				return
			}

			res <- utils.Result[uint64]{
				Val: writedImageHash,
				Err: nil,
			}
		}(i, image, results[i])
	}

	imageHashes := make([]utils.Result[uint64], 0, len(in.Images))
	for _, result := range results {
		imageHashes = append(imageHashes, <-result)
	}
	log.Trace().Msg("result readed from channels")

	// error check
	out.ImageErrors = make([]string, len(imageHashes))
	out.WritedIds = make([]string, len(imageHashes))
	errCounter := 0
	for i, imageHash := range imageHashes {
		if imageHash.Err != nil {
			out.ImageErrors[i] = imageHash.Err.Error()
			out.WritedIds[i] = "0"
			errCounter++
			if errCounter == 1 {
				out.Err = imageHash.Err.Error()
			}
			continue
		}

		out.ImageErrors[i] = "null"
		out.WritedIds[i] = strconv.FormatUint(imageHash.Val, 10)
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

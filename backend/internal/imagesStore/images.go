package imagesStore

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"

	"github.com/zpx64/supreme-octopus/internal/vars"

	"github.com/cespare/xxhash"
	"github.com/kolesa-team/go-webp/webp"
	webpDecoder "golang.org/x/image/webp"
)

const (
	LastBytesForHash = 256
)

func GetUniqHashFromImageFile(image []byte) uint64 {
	return xxhash.Sum64(image)
}

func PostImageToStore(imageBytes []byte, imageType string) (uint64, error) {
	var (
		imageReader = bytes.NewReader(imageBytes)
		imageHash   = GetUniqHashFromImageFile(imageBytes)

		imageDecoded image.Image
		err          error
	)

	storeForImages.imagesMapMutex.RLock()
	_, alreadyInStore := storeForImages.imagesMap[imageHash]
	storeForImages.imagesMapMutex.RUnlock()
	if alreadyInStore {
		return imageHash, nil
	}

	switch imageType {
	case "image/png":
		imageDecoded, err = png.Decode(imageReader)
		if err != nil {
			logger.Warn().Err(err).Msg("an error with image decoder")
			return 0, err
		}
	case "image/jpeg":
		imageDecoded, err = jpeg.Decode(imageReader)
		if err != nil {
			logger.Warn().Err(err).Msg("an error with image decoder")
			return 0, err
		}
	case "image/webp":
		imageDecoded, err = webpDecoder.Decode(imageReader)
		if err != nil {
			logger.Warn().Err(err).Msg("an error with image decoder")
			return 0, err
		}
	}

	// TODO: what the fuck is it
	if imageDecoded == nil {
		logger.Warn().Err(err).Msg("an error with image decoder: imageDecoded = nil")

		return 0, vars.ErrImageUpload
	}

	outFilePath := fmt.Sprintf(
		"%s/%s.webp",
		stateImagesDir,
		strconv.FormatUint(imageHash, 10),
	)

	outFile, err := os.Create(outFilePath)
	if err != nil {
		logger.Warn().Err(err).Msg("an error on output file creation")
		return 0, err
	}
	defer outFile.Close()

	err = webp.Encode(outFile, imageDecoded, webpEncoderOptions)
	if err != nil {
		logger.Warn().Err(err).Msg("an error with output file")
		return 0, err
	}

	storeForImages.imagesMapMutex.Lock()
	storeForImages.imagesMap[imageHash] = outFilePath
	storeForImages.imagesMapMutex.Unlock()

	return imageHash, nil
}

func GetImageFromStore(imageHash uint64) ([]byte, error) {
	storeForImages.imagesMapMutex.RLock()
	imageFilePath, ok := storeForImages.imagesMap[imageHash]
	storeForImages.imagesMapMutex.RUnlock()
	if !ok {
		return nil, vars.ErrImageNotFound
	}

	file, err := os.ReadFile(imageFilePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

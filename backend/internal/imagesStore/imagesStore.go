package imagesStore

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/zpx64/supreme-octopus/internal/vars"

	webpEncoder "github.com/kolesa-team/go-webp/encoder"
	"github.com/rs/zerolog"
)

type store struct {
	// TODO: maybe add support for user_id in info about image?
	imagesMap      map[uint64]string
	imagesMapMutex sync.RWMutex
	stateFile      *os.File
}

var (
	logger         zerolog.Logger
	stateFilePath  = vars.StatePath + "/images_store_state.json"
	stateImagesDir = vars.StatePath + "/images_store"

	storeForImages = store{}
	stopChan       = make(chan struct{})

	webpEncoderOptions *webpEncoder.Options

	inited = false
)

func SyncWithDiskState() error {
	storeForImages.imagesMapMutex.Lock()
	defer storeForImages.imagesMapMutex.Unlock()

	j, _ := json.Marshal(storeForImages.imagesMap)

	_, err := storeForImages.stateFile.Write(j)
	if err != nil {
		return err
	}

	return nil
}

func syncWithDiskStateInBackground() {
	ticker := time.NewTicker(
		time.Duration(vars.SyncStateTimeout) * time.Second,
	)
	for {
		select {
		case <-ticker.C:
			err := SyncWithDiskState()
			if err != nil {
				logger.Warn().Err(err).Send()
			}
		case <-stopChan:
			logger.Info().Msg("stoped syncing with disk")
			return
		}
	}
}

func Init(log *zerolog.Logger) error {
	if inited {
		return nil
	}

	logger = *log

	var (
		isStateFileExist = false
	)
	fileInfo, err := os.Stat(stateFilePath)
	if !errors.Is(err, os.ErrNotExist) {
		isStateFileExist = true
	}

	storeForImages.stateFile, err = os.OpenFile(
		stateFilePath,
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		0777,
	)
	if err != nil {
		logger.Warn().Err(err).Msg("an error with file creation")
		return err
	}

	storeForImages.imagesMap = make(map[uint64]string, vars.DefaultMapSize)
	if isStateFileExist && fileInfo.Size() > 2 {
		dec := json.NewDecoder(storeForImages.stateFile)

		err := dec.Decode(&storeForImages.imagesMap)
		if err != nil {
			logger.Warn().Err(err).Msg("an error with file reading")
			return err
		}
	}

	logger.Trace().Interface("imagesMap", storeForImages.imagesMap).Send()

	_, err = os.Stat(stateImagesDir)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(stateImagesDir, 0777)
		if err != nil {
			logger.Warn().Err(err).Msg("an error with dir creation")
			return err
		}
	}

	// TODO: test webp encoder lossy options
	webpEncoderOptions, err = webpEncoder.NewLossyEncoderOptions(
		webpEncoder.PresetDefault, 75,
	)
	if err != nil {
		return err
	}

	go syncWithDiskStateInBackground()
	inited = true
	return nil
}

// TODO: write normal crash handling
//
//	with shutdown flag
//	and rescan all image file on abnormal shutdown
//	arseniy dont ask me how to do it.
func Deinit() error {
	stopChan <- struct{}{}
	err := SyncWithDiskState()
	logger.Error().Err(err).Msg("an error with state file")
	storeForImages.stateFile.Close()
	storeForImages.imagesMapMutex.Lock()
	return err
}

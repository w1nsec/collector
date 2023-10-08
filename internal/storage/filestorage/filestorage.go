package filestorage

import (
	"bufio"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage/memstorage"
	"os"
)

type FileStorageInterface interface {
	memstorage.Storage
	Load() error
	SaveAll() error
	Close() error
	//Save(myMetrics metrics.MyMetrics) err
}

type FileStorage struct {
	filePath string
	file     *os.File
	memstorage.MemStorage
}

func (f FileStorage) Close() error {
	return f.file.Close()
}

func (f FileStorage) Load() error {
	// check that file exists, or already opened
	if f.file == nil {
		if _, err := os.Stat(f.filePath); err != nil {
			return err
		}
		file, err := os.OpenFile(f.filePath, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		f.file = file
	}

	sc := bufio.NewScanner(f.file)
	for sc.Scan() {
		var metric = &metrics.Metrics{}
		err := json.Unmarshal(sc.Bytes(), metric)
		if err != nil {
			return err
		}
		f.MemStorage.UpdateMetric(metric)

		metric = nil
	}
	if err := sc.Err(); err != nil {
		return err
	}
	return nil
}

func (f FileStorage) SaveAll() error {
	// check that file exists, or already opened
	if f.file == nil {
		if f.file == nil {
			if _, err := os.Stat(f.filePath); err != nil {
				return err
			}
			file, err := os.OpenFile(f.filePath, os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				return err
			}
			f.file = file
		}
	}

	mSlice := f.MemStorage.GetAllMetrics()
	for _, metric := range mSlice {
		encoder := json.NewEncoder(f.file)
		err := encoder.Encode(metric)
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}

	}

	return nil
}

package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage"
	"github.com/w1nsec/collector/internal/storage/memstorage"
)

type FileStorageInterface interface {
	//memstorage.Storage
	Load(ctx context.Context) error
	SaveAll(ctx context.Context) error
	Close(ctx context.Context) error
	//Save(myMetrics metrics.MyMetrics) err
}

type FileStorage struct {
	storage.Storage

	filePath string
	file     *os.File
	mutex    *sync.Mutex
}

func (f FileStorage) Close(context.Context) error {
	return f.file.Close()
}

func (f FileStorage) Load(ctx context.Context) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	// check that file exists, or already opened
	if f.file == nil {

		file, err := os.OpenFile(f.filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
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
		err = f.Storage.UpdateMetric(ctx, metric)
		if err != nil {
			log.Error().Err(err)
		}
		metric = nil
	}
	if err := sc.Err(); err != nil {
		return err
	}
	return nil
}

func (f FileStorage) SaveAll(ctx context.Context) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	// check that file exists, or already opened
	if f.file == nil {
		if _, err := os.Stat(f.filePath); err != nil {
			return err
		}
		file, err := os.OpenFile(f.filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		f.file = file
	}
	err := f.file.Truncate(0)
	_, err = f.file.Seek(0, 0)
	if err != nil {
		log.Error().Err(err).
			Msg("can't move to the beginning of file")
		return err
	}
	mSlice, err := f.Storage.GetAllMetrics(ctx)
	if err != nil {
		log.Error().Err(err).
			Msg("can't get all metrics value")
		return err
	}
	for _, metric := range mSlice {
		encoder := json.NewEncoder(f.file)
		err := encoder.Encode(metric)
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}
	}
	//defer f.file.Close()
	return nil
}

func NewFileStorage(path string, storage storage.Storage) (FileStorageInterface, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	f := &FileStorage{
		filePath: path,
		file:     file,
		Storage:  storage,
		mutex:    &sync.Mutex{},
	}
	if f.Storage == nil {
		f.Storage = memstorage.NewMemStorage()
	}
	return f, nil
}

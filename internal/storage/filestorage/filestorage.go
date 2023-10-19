package filestorage

import (
	"bufio"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage"
	"github.com/w1nsec/collector/internal/storage/memstorage"
	"os"
	"sync"
)

type FileStorageInterface interface {
	//memstorage.Storage
	Load() error
	SaveAll() error
	Close() error
	//Save(myMetrics metrics.MyMetrics) err
}

type FileStorage struct {
	storage.Storage

	filePath string
	file     *os.File
	mutex    *sync.Mutex
}

func (f FileStorage) Close() error {
	return f.file.Close()
}

func (f FileStorage) Load() error {
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
		f.Storage.UpdateMetric(metric)

		metric = nil
	}
	if err := sc.Err(); err != nil {
		return err
	}
	return nil
}

func (f FileStorage) SaveAll() error {
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
	f.file.Truncate(0)
	f.file.Seek(0, 0)
	mSlice, _ := f.Storage.GetAllMetrics()
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

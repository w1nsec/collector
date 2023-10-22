package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config/server"
	"github.com/w1nsec/collector/internal/logger"
	"github.com/w1nsec/collector/internal/storage"
	"github.com/w1nsec/collector/internal/storage/filestorage"
	"sync"
	"time"
)

type MetricService struct {
	// context
	ctx context.Context

	// storage
	storage.Storage
	filestorage.FileStorageInterface

	StoreInterval time.Duration
	Restore       bool
}

func (service *MetricService) CheckStorage() error {
	return service.Storage.CheckStorage()
}

func (service *MetricService) Close() error {
	err := service.FileStorageInterface.Close()
	if err != nil {
		return err
	}
	return service.Storage.Close()
}

func (service *MetricService) SetupLogger(level string) error {
	return logger.Initialize(level)
}

func NewService(args server.Args, store storage.Storage,
	fstore filestorage.FileStorageInterface) (*MetricService, error) {

	// create service
	service := &MetricService{
		ctx:                  context.Background(),
		Storage:              store,
		FileStorageInterface: fstore,
		Restore:              args.Restore,
		StoreInterval:        time.Duration(args.StoreInterval) * time.Second,
	}
	err := service.SetupLogger(args.LogLevel)
	if err != nil {
		return nil, err
	}
	return service, err
}

func (service *MetricService) BackupLoop(wg *sync.WaitGroup, storeInterval time.Duration) error {
	timer := time.NewTicker(storeInterval)
	for t := range timer.C {
		err := service.FileStorageInterface.SaveAll()
		if err != nil {
			log.Info().
				Str("time", t.Format(time.DateTime)).
				Msg("DB saved")
			continue
		}

	}
	wg.Done()
	return nil
}

func (service *MetricService) Setup(wg *sync.WaitGroup) error {
	if service.Restore {
		if service.FileStorageInterface != nil {
			service.FileStorageInterface.Load()
			wg.Add(1)
			go service.BackupLoop(wg, service.StoreInterval)
		}
		return service.Storage.Init()
	}
	return nil
}

package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config/server"
	"github.com/w1nsec/collector/internal/logger"
	"github.com/w1nsec/collector/internal/storage"
	"github.com/w1nsec/collector/internal/storage/filestorage"
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

	// iter14, hmac signing
	Secret string
}

func (service *MetricService) CheckStorage() error {
	return service.Storage.CheckStorage()
}

func (service *MetricService) Close(ctx context.Context) error {
	if service.FileStorageInterface != nil {
		err := service.FileStorageInterface.Close(ctx)
		if err != nil {
			defer service.Storage.Close(ctx)
			return err
		}

	}

	return service.Storage.Close(ctx)
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
		Secret:               args.Key,
	}
	err := service.SetupLogger(args.LogLevel)
	if err != nil {
		return nil, err
	}
	return service, err
}

func (service *MetricService) BackupLoop(ctx context.Context, storeInterval time.Duration) error {
	timer := time.NewTicker(storeInterval)
	for {
		select {
		case t := <-timer.C:
			err := service.FileStorageInterface.SaveAll(ctx)
			if err != nil {
				log.Info().
					Str("time", t.Format(time.DateTime)).
					Msg("DB saved")
				continue
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (service *MetricService) Setup(ctx context.Context) error {
	if service.Restore {
		if service.FileStorageInterface != nil {
			service.FileStorageInterface.Load(ctx)
			go service.BackupLoop(ctx, service.StoreInterval)
		}
		return service.Storage.Init()
	}
	return nil
}

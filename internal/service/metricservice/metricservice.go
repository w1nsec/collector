package metricservice

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config/server"
	"github.com/w1nsec/collector/internal/logger"
	"github.com/w1nsec/collector/internal/server/http"
	"github.com/w1nsec/collector/internal/service"
	"github.com/w1nsec/collector/internal/storage"
	"github.com/w1nsec/collector/internal/storage/dbstorage"
	"github.com/w1nsec/collector/internal/storage/filestorage"
	"github.com/w1nsec/collector/internal/storage/memstorage"
	"sync"
	"time"
)

type MetricService struct {
	// context
	ctx context.Context

	// web server
	server http.Server

	// storage
	storage.Storage
	filestorage.FileStorageInterface
	db dbstorage.DBStorage

	StoreInterval time.Duration
	Restore       bool
}

func (service *MetricService) CheckDB() error {
	if service.db == nil {
		return fmt.Errorf("db not used")
	}
	return service.db.CheckConnection()
}

func (service *MetricService) Close() error {
	err := service.FileStorageInterface.Close()
	if err != nil {
		return err
	}
	return service.db.Close()
}

func (service *MetricService) SetupLogger(level string) error {
	return logger.Initialize(level)
}

func NewService(args server.Args) (service.Service, error) {
	// initialise storages
	var (
		store   storage.Storage
		fstore  filestorage.FileStorageInterface
		dbstore dbstorage.DBStorage
	)

	// create service
	service := &MetricService{
		ctx:           context.Background(),
		StoreInterval: time.Duration(args.StoreInterval) * time.Second,
		Restore:       args.Restore,
	}

	if args.DatabaseURL != "" {
		dbstore = dbstorage.NewPostgresStorage(args.DatabaseURL)
	} else {
		store = memstorage.NewMemStorage()
		var err error
		fstore, err = filestorage.NewFileStorage(args.StoragePath, store)
		if err != nil {
			return nil, err
		}
	}

	// initialise server
	server, err := http.NewServerForService(args, service)
	if err != nil {
		return nil, err
	}
	mux := http.NewRouter(service)
	server.AddMux(mux)

	// add server to service
	service.server = server

	// init db storage
	//service.db = args.DatabaseURL)
	if dbstore == nil {
		service.Storage = store
		service.FileStorageInterface = fstore
		service.db = nil
	} else {
		service.Storage = dbstore
		service.db = dbstore
	}

	return service, nil
}

func (service MetricService) Start() error {
	// restore DB
	wg := &sync.WaitGroup{}
	if service.Restore {
		if service.FileStorageInterface != nil {
			service.FileStorageInterface.Load()
			wg.Add(1)
			go service.BackupLoop(wg)
		}
		if service.db != nil {
			err := service.db.CreateTables()
			if err != nil {
				return err
			}
		}
	}

	// start server
	err := service.server.Start()
	wg.Wait()
	return err
}

func (service *MetricService) BackupLoop(wg *sync.WaitGroup) error {
	timer := time.NewTicker(service.StoreInterval)
	for t := range timer.C {
		log.Info().
			Str("time", t.Format(time.DateTime)).
			Msg("DB saved")

		service.FileStorageInterface.SaveAll()
	}

	wg.Done()
	return nil
}

func (service MetricService) Stop() error {

	log.Error().
		Err(service.FileStorageInterface.SaveAll()).
		Msg("fs-storage saved")

	log.Error().
		Err(service.FileStorageInterface.Close()).
		Msg("fs-storage closed")
	defer service.db.Close()
	return service.server.Stop()
}

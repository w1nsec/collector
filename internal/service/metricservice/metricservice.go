package metricservice

import (
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config"
	"github.com/w1nsec/collector/internal/logger"
	"github.com/w1nsec/collector/internal/server/http"
	"github.com/w1nsec/collector/internal/service"
	"github.com/w1nsec/collector/internal/storage/dbstorage"
	"github.com/w1nsec/collector/internal/storage/filestorage"
	"github.com/w1nsec/collector/internal/storage/memstorage"
	"sync"
	"time"
)

type MetricService struct {
	server http.Server

	// storage
	memstorage.Storage
	filestorage.FileStorageInterface
	dbstorage.DBStorage

	StoreInterval time.Duration
	Restore       bool
}

func (service *MetricService) CheckDB() error {
	return service.CheckConnection()
}

func (service *MetricService) Close() error {
	err := service.FileStorageInterface.Close()
	if err != nil {
		return err
	}
	return service.DBStorage.Close()
}

func (service *MetricService) SetupLogger(level string) error {
	return logger.Initialize(level)
}

func NewService(args config.Args) (service.Service, error) {
	// initialise storages
	store := memstorage.NewMemStorage()
	fstore, err := filestorage.NewFileStorage(args.StoragePath, store)
	if err != nil {
		return nil, err
	}

	service := &MetricService{
		server:               nil,
		Storage:              store,
		FileStorageInterface: fstore,
		StoreInterval:        time.Duration(args.StoreInterval) * time.Second,
		Restore:              args.Restore,
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
	service.DBStorage = dbstorage.NewPostgresStorage(args.DatabaseURL)

	return service, nil
}

func (service MetricService) Start() error {
	// restore DB
	if service.Restore {
		service.FileStorageInterface.Load()
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go service.BackupLoop(wg)

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
	defer service.DBStorage.Close()
	return service.server.Stop()
}

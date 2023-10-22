package app

import (
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config/server"
	"github.com/w1nsec/collector/internal/server/http"
	"github.com/w1nsec/collector/internal/service"
	"github.com/w1nsec/collector/internal/storage"
	"github.com/w1nsec/collector/internal/storage/dbstorage"
	"github.com/w1nsec/collector/internal/storage/filestorage"
	"github.com/w1nsec/collector/internal/storage/memstorage"
	"sync"
)

type appServer struct {
	service *service.MetricService

	// now it is internal/server/http
	server http.Server // interfaces

}

func NewAppServer() (*appServer, error) {
	var (
		store  storage.Storage
		fstore filestorage.FileStorageInterface
		args   server.Args
	)

	server.ServerArgsParse(&args)
	log.Info().
		Str("addr", args.Addr).
		Str("log", args.LogLevel).Send()

	if args.DatabaseURL != "" {
		store = dbstorage.NewStorage(args.DatabaseURL)
	} else {
		store = memstorage.NewMemStorage()
		var err error
		fstore, err = filestorage.NewFileStorage(args.StoragePath, store)
		if err != nil {
			return nil, err
		}
	}

	// initialise service
	service, err := service.NewService(args, store, fstore)
	if err != nil {
		return nil, err
	}

	// initialise transport server
	server, err := http.NewServerForService(args, service)
	if err != nil {
		return nil, err
	}

	app := &appServer{
		service: service,
		server:  server,
	}

	return app, nil
}

func (app appServer) Run() error {
	// initialise storages

	// restore DB
	wg := &sync.WaitGroup{}
	err := app.service.Setup(wg)
	if err != nil {
		return err
	}
	// start server
	err = app.server.Start()
	wg.Wait()
	return err
}

func (app *appServer) Stop() error {

	log.Error().
		Err(app.service.FileStorageInterface.SaveAll()).
		Msg("fs-storage saving status")
	log.Error().
		Err(app.service.FileStorageInterface.Close()).
		Msg("fs-storage closing status")
	defer app.service.Close()
	return app.server.Stop()
}

func (app appServer) CheckStorage() bool {
	return app.service.Storage != nil
}

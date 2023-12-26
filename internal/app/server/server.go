package server

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config/server"
	"github.com/w1nsec/collector/internal/service"
	"github.com/w1nsec/collector/internal/storage"
	"github.com/w1nsec/collector/internal/storage/dbstorage"
	"github.com/w1nsec/collector/internal/storage/filestorage"
	"github.com/w1nsec/collector/internal/storage/memstorage"
	"github.com/w1nsec/collector/internal/transport/http"
)

type appServer struct {
	service *service.MetricService

	// now it is internal/transport/http
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

	// initialise transport transport
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

func (app appServer) Run(ctx context.Context) error {
	// initialise storages

	// restore DB
	go app.service.Setup(ctx)

	// start transport
	go func() {
		err := app.server.Start()
		if err != nil {
			log.Error().Err(err).
				Msg("Error starting transport")
		}
	}()

	// get signal to exit
	<-ctx.Done()

	log.Info().Str("signal", "Ctlr+C").Send()
	shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	longShutdown := make(chan struct{}, 1)
	go func() {
		time.Sleep(10 * time.Second)
		longShutdown <- struct{}{}
	}()

	// if shutdown too long
	select {
	case <-shutdownCtx.Done():
		err := fmt.Errorf("transport shutdown: %v", ctx.Err())
		log.Error().Err(err).Send()
		err = app.Stop(ctx)
		return err
	case <-longShutdown:
		err := fmt.Errorf("force finishing")
		log.Error().Err(err).Send()
		return err
	}
}

func (app *appServer) Stop(ctx context.Context) error {
	if app.service.FileStorageInterface != nil {
		log.Error().
			Err(app.service.FileStorageInterface.SaveAll(ctx)).
			Msg("fs-storage saving status")
		log.Error().
			Err(app.service.FileStorageInterface.Close(ctx)).
			Msg("fs-storage closing status")
	}

	err := app.service.Close(ctx)
	if err != nil {
		log.Error().Err(err).Send()
	}

	return app.server.Stop()
}

func (app appServer) CheckStorage() bool {
	return app.service.Storage != nil
}

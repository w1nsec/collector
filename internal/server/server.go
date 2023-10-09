package server

import (
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config"
	"github.com/w1nsec/collector/internal/handlers"
	"github.com/w1nsec/collector/internal/logger"
	"github.com/w1nsec/collector/internal/storage/filestorage"
	"github.com/w1nsec/collector/internal/storage/memstorage"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server interface {
	InitLogger(level string) error
	Start() error
}

type MetricServer struct {
	addr *net.TCPAddr
	//Store *memstorage.MemStorage
	Store         memstorage.Storage
	FileStore     filestorage.FileStorageInterface
	StoreInterval time.Duration
	Restore       bool

	//mux   *http.ServeMux
	http.Server
}

func (srv *MetricServer) InitLogger(level string) error {
	return logger.Initialize(level)
}

// NewServer is just a wrapper for NewMetricServerWithParams
// with default params, return interface for server
// func NewServer(addr string, loggerLevel string) (Server, error) {
func NewServer(args config.Args) (Server, error) {
	store := memstorage.NewMemStorage()
	mux := handlers.NewRouter(store)
	return NewMetricServerWithParams(args.Addr, store, mux, args.LogLevel, args.Restore, args.StoragePath, args.StoreInterval)
}

func NewMetricServerWithParams(addr string,
	store memstorage.Storage,
	mux http.Handler, loggerLevel string, restore bool,
	storePath string, storInterval uint64) (*MetricServer, error) {

	netAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	if store == nil {
		store = memstorage.NewMemStorage()
	}

	fstore, err := filestorage.NewFileStorage(storePath, store)
	if err != nil {
		return nil, err
	}

	if mux == nil {
		//mux = http.NewServeMux()
		mux = handlers.NewRouter(store)

	}

	srv := MetricServer{
		addr:  netAddr,
		Store: store,
		Server: http.Server{
			Addr:    netAddr.String(),
			Handler: mux,
		},
		FileStore:     fstore,
		StoreInterval: time.Duration(storInterval) * time.Second,
	}

	err = srv.InitLogger(loggerLevel)
	if err != nil {
		return nil, err
	}
	return &srv, nil
}

func (srv *MetricServer) AddMux(mux *http.ServeMux) {
	srv.Server.Handler = mux
	//srv.mux = mux
}

func (srv *MetricServer) Close() {
	srv.FileStore.SaveAll()
	srv.FileStore.Close()
	srv.Server.Close()
}

func (srv *MetricServer) Start() error {
	// restore DB
	if srv.Restore {
		srv.FileStore.Load()
	}

	// Start server
	log.Info().Msgf("[+] Started on: %s", srv.Addr)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go srv.BackupLoop(wg)

	defer srv.Close()
	err := srv.ListenAndServe()
	wg.Wait()
	return err
}

func (srv *MetricServer) BackupLoop(wg *sync.WaitGroup) error {
	timer := time.NewTicker(srv.StoreInterval)
	for t := range timer.C {
		log.Info().
			Str("time", t.Format(time.DateTime)).
			Msg("DB saved")

		srv.FileStore.SaveAll()
	}

	wg.Done()
	return nil
}

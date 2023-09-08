package server

import (
	"github.com/w1nsec/collector/internal/memstorage"
	"net"
	"net/http"
)

type MetricServer struct {
	addr *net.TCPAddr
	//Store *memstorage.MemStorage
	Store memstorage.Storage
	//mux   *http.ServeMux
	http.Server
}

func NewMetricServer(addr string, store memstorage.Storage, mux *http.ServeMux) (*MetricServer, error) {
	netAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	if store == nil {
		store = memstorage.NewMemStorage()
	}
	if mux == nil {
		mux = http.NewServeMux()
	}

	srv := MetricServer{
		addr:  netAddr,
		Store: store,
		Server: http.Server{
			Addr:    netAddr.String(),
			Handler: mux,
		},
	}
	return &srv, nil
}

func (srv *MetricServer) AddMux(mux *http.ServeMux) {
	srv.Server.Handler = mux
	//srv.mux = mux
}

func (srv *MetricServer) Start() error {
	return srv.ListenAndServe()
}

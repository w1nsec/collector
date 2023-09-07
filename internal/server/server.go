package server

import (
	"github.com/w1nsec/collector/internal/memstorage"
	"net"
	"net/http"
)

type MetricServer struct {
	addr  *net.TCPAddr
	Store *memstorage.MemStorage
	mux   *http.ServeMux
}

func NewMetricServer(addr string, store *memstorage.MemStorage, mux *http.ServeMux) (*MetricServer, error) {
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
		mux:   mux,
	}
	return &srv, nil
}

func (srv *MetricServer) AddMux(mux *http.ServeMux) {
	srv.mux = mux
}

func (srv *MetricServer) Start() error {
	return http.ListenAndServe(srv.addr.String(), srv.mux)
}

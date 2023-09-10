package server

import (
	"fmt"
	"github.com/w1nsec/collector/internal/handlers"
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

func NewMetricServer(addr string) (*MetricServer, error) {
	store := memstorage.NewMemStorage()
	mux := handlers.NewRouter(store)
	return NewMetricServerWithParams(addr, store, mux)
}
func NewMetricServerWithParams(addr string, store memstorage.Storage, mux http.Handler) (*MetricServer, error) {
	netAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	if store == nil {
		store = memstorage.NewMemStorage()
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
	}
	return &srv, nil
}

func (srv *MetricServer) AddMux(mux *http.ServeMux) {
	srv.Server.Handler = mux
	//srv.mux = mux
}

func (srv *MetricServer) Start() error {
	fmt.Println("[+] Started on:", srv.Addr)
	return srv.ListenAndServe()
}

package http

import (
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config/server"
	"github.com/w1nsec/collector/internal/service"
)

type Server interface {
	Start() error
	Stop() error

	AddMux(mux http.Handler)
}

type MetricServer struct {
	addr *net.TCPAddr

	//mux   *http.ServeMux
	http.Server
}

func NewServerForService(args server.Args, service *service.MetricService) (Server, error) {

	mux := NewRouter(service)
	return NewMetricServerWithParams(args.Addr, mux)
}

func NewMetricServerWithParams(addr string, mux http.Handler) (*MetricServer, error) {

	netAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	srv := MetricServer{
		addr: netAddr,
		Server: http.Server{
			Addr:    netAddr.String(),
			Handler: mux,
		},
	}

	return &srv, nil
}

func (srv *MetricServer) AddMux(mux http.Handler) {
	srv.Server.Handler = mux
	//srv.mux = mux
}

func (srv *MetricServer) Stop() error {
	return srv.Server.Close()
}

func (srv *MetricServer) Start() error {
	// Start transport
	log.Info().Msgf("[+] Started on: %s", srv.Addr)
	return srv.ListenAndServe()
}

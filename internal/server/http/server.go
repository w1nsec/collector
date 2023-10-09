package http

import (
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config"
	"net"
	"net/http"
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

// NewServer is just a wrapper for NewMetricServerWithParams
// with default params, return interface for server
// func NewServer(addr string, loggerLevel string) (Server, error) {
func NewServer(args config.Args) (Server, error) {

	//mux := handlers.NewRouter(service)
	return NewMetricServerWithParams(args.Addr)
}

func NewMetricServerWithParams(addr string) (*MetricServer, error) {

	netAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	srv := MetricServer{
		addr: netAddr,
		Server: http.Server{
			Addr: netAddr.String(),
			//Handler: mux,
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
	// Start server
	log.Info().Msgf("[+] Started on: %s", srv.Addr)
	return srv.ListenAndServe()
}

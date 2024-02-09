package http

import (
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config/server"
	"github.com/w1nsec/collector/internal/service"
)

// @Title MetricServer Endpoints
// @Description Service for storing metrics.
// @Version 1.0
// @BasePath /
// @Host localhost:8080

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

func NewServerForService(args *server.Args, service *service.MetricService) (Server, error) {

	_, cidr, _ := net.ParseCIDR(args.CIDR)
	mux := NewRouter(service, cidr)

	return NewMetricServerWithParams(args.Addr, mux, args.CIDR)
}

func NewMetricServerWithParams(addr string, mux http.Handler, netArg string) (*MetricServer, error) {
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

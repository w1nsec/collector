package grpc

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/service"
	pb "github.com/w1nsec/collector/internal/transport/grpc/proto"
	conv "github.com/w1nsec/collector/internal/utils/proto_conv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MetricsGRPC struct {
	addr   *net.TCPAddr
	srvRPC *grpc.Server
	srvMet *MetricsServer
}

func NewMetricsGRPC(addr string, mSRV *MetricsServer) (srv *MetricsGRPC, err error) {
	netAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()
	srv = &MetricsGRPC{
		addr:   netAddr,
		srvRPC: s,
		srvMet: mSRV,
	}

	return srv, nil
}

func (s *MetricsGRPC) Stop() error {
	s.srvRPC.GracefulStop()
	log.Info().Msgf("[+] Server stopped")
	return nil

}

func (s *MetricsGRPC) Start() error {
	l, err := net.Listen("tcp", s.addr.String())
	if err != nil {
		log.Error().
			Err(err).Msgf("can't start server")
		return err
	}

	pb.RegisterMetricsSvcServer(s.srvRPC, s.srvMet)

	log.Info().Msgf("[+] Started on: %s", s.addr.String())
	return s.srvRPC.Serve(l)
}

type MetricsServer struct {
	pb.UnimplementedMetricsSvcServer

	service *service.MetricService
}

func NewMetricsServer(service *service.MetricService) (srv *MetricsServer, err error) {
	return &MetricsServer{service: service}, nil
}

func (s *MetricsServer) StoreMetrics(ctx context.Context,
	inMet *pb.Metrics) (resp *pb.StoreResponse, err error) {
	resp = &pb.StoreResponse{}

	ms := conv.ConvertProto2Metrics(inMet)
	err = s.service.UpdateMetrics(ctx, ms)
	if err != nil && err.Error() != "" {
		return nil, status.Errorf(codes.Internal, "can't add/update metrics: %v", err)
	}

	return resp, nil
}

func (s *MetricsServer) ListMetrics(ctx context.Context,
	reqLimit *pb.ListMetricsReq) (allM *pb.Metrics, err error) {

	allM = new(pb.Metrics)

	ms, err := s.service.GetAllMetrics(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't get all metrics: %v", err)
	}

	return conv.ConvertMetrics2Proto(ms), nil
}

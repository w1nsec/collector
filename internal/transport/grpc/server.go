package grpc

import (
	"context"

	"github.com/w1nsec/collector/internal/service"
	pb "github.com/w1nsec/collector/internal/transport/grpc/proto"
	conv "github.com/w1nsec/collector/internal/utils/proto_conv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MetricsServer struct {
	pb.UnimplementedMetricsSvcServer
	//addr    *net.Addr
	service *service.MetricService
}

func NewMetricsServer(service *service.MetricService) *MetricsServer {
	return &MetricsServer{service: service}
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

package conversion

import (
	"github.com/w1nsec/collector/internal/metrics"
	pb "github.com/w1nsec/collector/internal/transport/grpc/proto"
)

func ConvertMetrics2Proto(m []*metrics.Metrics) *pb.Metrics {

	mSlice := make([]*pb.Metric, len(m))
	for i, val := range m {
		switch val.MType {
		case metrics.Counter:
			mSlice[i] = &pb.Metric{
				Id:    val.ID,
				Mtype: pb.Metric_COUNTER,
				Delta: *val.Delta,
			}
		case metrics.Gauge:
			mSlice[i] = &pb.Metric{
				Id:    val.ID,
				Mtype: pb.Metric_GAUGE,
				Value: *val.Value,
			}
		}
	}
	return &pb.Metrics{
		Metrics: mSlice,
	}
}

func ConvertProto2Metrics(m *pb.Metrics) []*metrics.Metrics {
	sl := make([]*metrics.Metrics, len(m.Metrics))
	for ind, val := range m.Metrics {
		sl[ind] = new(metrics.Metrics)
		switch val.Mtype {
		case pb.Metric_GAUGE:
			sl[ind].ID = val.Id
			sl[ind].MType = metrics.Gauge
			sl[ind].Value = &val.Value
		case pb.Metric_COUNTER:
			sl[ind].ID = val.Id
			sl[ind].MType = metrics.Counter
			sl[ind].Delta = &val.Delta
		}
	}

	return sl
}

package agent

import (
	"context"
	"fmt"
	"time"

	pb "github.com/w1nsec/collector/internal/transport/grpc/proto"
	conversion "github.com/w1nsec/collector/internal/utils/proto_conv"

	"github.com/w1nsec/collector/internal/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (agent Agent) SendRPCBatch(job []*metrics.Metrics) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	c, err := grpc.Dial(agent.addr.String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("can't connect to gRPC server")
	}

	cli := pb.NewMetricsSvcClient(c)
	ms := conversion.ConvertMetrics2Proto(job)
	_, err = cli.StoreMetrics(ctx, ms)
	if err != nil {
		return fmt.Errorf("can't store sent metrics")
	}
	return nil
}

package grpc

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/w1nsec/collector/internal/config/server"
	"github.com/w1nsec/collector/internal/service"
	"github.com/w1nsec/collector/internal/storage/memstorage"
	pb "github.com/w1nsec/collector/internal/transport/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestNewMetricsServer(t *testing.T) {
	mem := memstorage.NewMemStorage()
	svc, err := service.NewService(&server.Args{}, mem, nil)
	require.NoError(t, err)

	srv, err := NewMetricsServer(svc)
	require.NoError(t, err)
	require.NotNil(t, srv)
	require.Equal(t, mem, srv.service.Storage)

}

func TestNewMetricsGRPC(t *testing.T) {
	mem := memstorage.NewMemStorage()
	svc, err := service.NewService(&server.Args{}, mem, nil)
	require.NoError(t, err)

	srv, err := NewMetricsServer(svc)
	require.NoError(t, err)
	require.NotNil(t, srv)
	require.Equal(t, mem, srv.service.Storage)

	addr := "127.0.0.1:9999"
	srvRPC, err := NewMetricsGRPC(addr, srv)
	require.NoError(t, err)
	require.NotNil(t, srvRPC)
	require.Equal(t, srv, srvRPC.srvMet)
	require.Equal(t, addr, srvRPC.addr.String())

	addrWrong := "127.0.0.1:34773483"
	_, err = NewMetricsGRPC(addrWrong, srv)
	require.Error(t, err)
}

func TestMetricsServer(t *testing.T) {
	ms := []*pb.Metric{
		{Id: "test1_gauge", Mtype: pb.Metric_GAUGE, Value: 234.2343},
		{Id: "test2_counter", Mtype: pb.Metric_COUNTER, Delta: 234},
	}
	all := pb.Metrics{Metrics: ms}
	fmt.Println(all.Metrics)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	start := int64(30000)
	r, _ := rand.Int(rand.Reader, big.NewInt(65535-start))
	port := start + r.Int64()
	addr := fmt.Sprintf("localhost:%d", port)

	mem := memstorage.NewMemStorage()
	svc, err := service.NewService(&server.Args{}, mem, nil)
	require.NoError(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	srvM, err := NewMetricsServer(svc)
	require.NoError(t, err)

	srv, err := NewMetricsGRPC(addr, srvM)
	require.NoError(t, err)

	go srv.Start()
	defer srv.Close()

	time.Sleep(time.Second * 1)

	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		wg.Done()
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewMetricsSvcClient(conn)

	metrics, err := c.ListMetrics(ctx, &pb.ListMetricsReq{})
	if err != nil {
		wg.Done()
		require.NoError(t, err)
	}

	fmt.Println("Len:", len(metrics.Metrics))
	for i, v := range metrics.Metrics {
		fmt.Printf("%d: ID: %s | Val: %f | Del: %d\n", i, v.Id, v.Value, v.Delta)
	}

	fmt.Println("---- Storing metrics ----")
	_, err = c.StoreMetrics(ctx, &all)
	if err != nil {
		wg.Done()
		require.NoError(t, err)
	}

	metrics, err = c.ListMetrics(ctx, &pb.ListMetricsReq{})
	if err != nil {
		wg.Done()
		log.Fatal(err)

	}

	fmt.Println("Len:", len(metrics.Metrics))
	for i, v := range metrics.Metrics {
		fmt.Printf("%d: ID: %s | Val: %f | Del: %d\n", i, v.Id, v.Value, v.Delta)
	}
	require.Equal(t, len(all.Metrics), len(metrics.Metrics))

	wg.Done()
	wg.Wait()
}

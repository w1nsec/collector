// Package agent contains Agent struct
// that provide functionality for collecting and sending
// metrics to server part
package agent

import (
	"context"
	"net"
	"net/http"
	"time"

	config "github.com/w1nsec/collector/internal/config/agent"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage"
	"github.com/w1nsec/collector/internal/storage/memstorage"
)

// Params for sleeping agent if it receives too many errors
var (
	timeout = 10 * time.Second
	//retryStep     = uint(2) // 2 seconds
	maxRetryCount = uint(5)
)

// Agent struct, that contains Storage and other config options for running agent
type Agent struct {
	addr         net.Addr
	metricsPoint string

	store          storage.Storage
	pollInterval   time.Duration
	reportInterval time.Duration
	compression    bool
	httpClient     *http.Client

	// increment 13
	retryCount uint
	sleepCount uint
	sleepCh    []chan struct{}

	// param needs for signing send body (increment 14)
	secret string

	// increment 15
	rateLimit     int
	errorsCh      chan error
	successReport chan struct{}
}

// NewAgent is constructor for Agent struct
func NewAgent(args *config.Args) (*Agent, error) {
	netAddr, err := net.ResolveTCPAddr("tcp", args.Addr)
	if err != nil {
		return nil, err
	}

	sleepCh := make([]chan struct{}, args.Rate)
	for id := range sleepCh {
		sleepCh[id] = make(chan struct{})
	}

	agent := &Agent{
		addr:           netAddr,
		metricsPoint:   "update",
		pollInterval:   time.Duration(args.PollInterval) * time.Second,
		reportInterval: time.Duration(args.ReportInterval) * time.Second,

		store:      memstorage.NewMemStorage(),
		httpClient: &http.Client{Timeout: timeout},

		// increment 13
		retryCount: maxRetryCount,
		sleepCount: 0,
		sleepCh:    sleepCh,

		// increment 14
		secret: args.Key,

		// increment 15
		rateLimit:     args.Rate,
		errorsCh:      make(chan error, args.Rate),
		successReport: make(chan struct{}, args.Rate),
	}

	agent.compression = true
	return agent, nil
}

// Start run Agent functionality (collect and send metrics)
func (agent Agent) Start(ctx context.Context) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	metricsChannel := make(chan []*metrics.Metrics)

	// fill agent store
	go func() {
		agent.generator(ctx, metricsChannel)
	}()

	// create workers pull
	go func() {
		agent.limiter(ctx, metricsChannel)
	}()

	// validate localerrors count
	go func() {
		agent.validateErrors(ctx)
	}()

	// waiting, until goroutines not done
	<-ctx.Done()

	return nil
}

// Close channels
func (agent Agent) Close() {
	close(agent.errorsCh)
	close(agent.successReport)
	agent.store.Close(context.TODO())
}

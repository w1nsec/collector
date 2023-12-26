package agent

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	config "github.com/w1nsec/collector/internal/config/agent"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage"
	"github.com/w1nsec/collector/internal/storage/memstorage"
)

var usedMemStats = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

var (
	timeout = 10 * time.Second
	//retryStep     = uint(2) // 2 seconds
	maxRetryCount = uint(5)
)

type Agent struct {
	addr         net.Addr
	metricsPoint string

	// unused field
	metrics map[string]metrics.MyMetrics
	////

	store          storage.Storage
	pollInterval   time.Duration
	reportInterval time.Duration
	compression    bool
	httpClient     *http.Client

	// increment 13
	retryCount uint
	sleepCount uint
	sleepCh    []chan struct{}

	// increment 14
	secret string

	// increment 15
	rateLimit     int
	errorsCh      chan error
	successReport chan struct{}
}

func NewAgent(args config.Args) (*Agent, error) {
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
		metrics:        make(map[string]metrics.MyMetrics),
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

func (agent Agent) StartOLD(ctx context.Context) error {

	var (
		curErrCount = uint(0)
	)

	pollTicker := time.NewTicker(agent.pollInterval)
	reportTicker := time.NewTicker(agent.reportInterval)
	for {
		select {
		case t1 := <-pollTicker.C:
			fmt.Println("Receiving:", t1.Format(time.TimeOnly))
			agent.CollectMetrics(ctx)
		case t2 := <-reportTicker.C:
			fmt.Println("- Sending:", t2.Format(time.TimeOnly))

			// send all metrics together
			err := agent.SendAllMetricsJSON()
			if err != nil {
				log.Debug().
					Msgf("%v error, while send all metrics together", err)
				curErrCount += 1
				if errors.Is(err, syscall.ECONNREFUSED) {
					// agent can't connect to transport, let's try again
					sleep := 1 * time.Second
					for i := uint(0); i < agent.retryCount; i++ {
						time.Sleep(sleep)
						err := agent.SendAllMetricsJSON()
						// continue if success
						if err == nil {
							break
						}
						// update sleep time
						sleep = time.Duration(sleep.Seconds() + float64(agent.sleepCount))
					}
				}
				if curErrCount > agent.retryCount {
					return err
				}
			}
			curErrCount = 0
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
			defer cancel()
			return agent.store.Close(shutdownCtx)
		}
	}
}

func (agent Agent) Start(ctx context.Context) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// TODO maybe, metricsChannel capacity should be agent.rateLimit ??
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

func (agent Agent) Close() {
	close(agent.errorsCh)
	close(agent.successReport)
	agent.store.Close(context.TODO())
}

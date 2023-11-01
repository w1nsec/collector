package agent

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	config "github.com/w1nsec/collector/internal/config/agent"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage"
	"github.com/w1nsec/collector/internal/storage/memstorage"
	"net"
	"net/http"
	"syscall"
	"time"
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
	timeout       = 10 * time.Second
	retryStep     = uint(2) // 2 seconds
	maxRetryCount = uint(3)
)

type Agent struct {
	addr           net.Addr
	metricsPoint   string
	metrics        map[string]metrics.MyMetrics
	store          storage.Storage
	pollInterval   time.Duration
	reportInterval time.Duration
	compression    bool
	httpClient     *http.Client

	// increment 13
	retryCount uint
	retryStep  uint

	// increment 14
	secret string
}

func NewAgent(args config.Args) (*Agent, error) {
	netAddr, err := net.ResolveTCPAddr("tcp", args.Addr)
	if err != nil {
		return nil, err
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
		retryStep:  retryStep,

		// increment 14
		secret: args.Key,
	}

	agent.compression = true
	return agent, nil
}

func (agent Agent) Start() error {

	var (
		curErrCount = uint(0)
	)

	// Receive and send for the first time
	//fmt.Println("Receiving:", time.Now().Format(time.TimeOnly))
	//agent.GetMetrics()
	//fmt.Println("- Sending:", time.Now().Format(time.TimeOnly))
	//agent.SendMetrics()
	//err := agent.SendMetricsJSON()
	//if err != nil {
	//	return err
	//}

	pollTicker := time.NewTicker(agent.pollInterval)
	reportTicker := time.NewTicker(agent.reportInterval)
	for {
		select {
		case t1 := <-pollTicker.C:
			fmt.Println("Receiving:", t1.Format(time.TimeOnly))
			//agent.GetMetrics()
			agent.CollectMetrics()
		case t2 := <-reportTicker.C:
			fmt.Println("- Sending:", t2.Format(time.TimeOnly))

			// send all metrics together
			err := agent.SendAllMetricsJSON()
			if err != nil {
				log.Debug().
					Msgf("%v error, while send all metrics together", err)
				curErrCount += 1
				if errors.Is(err, syscall.ECONNREFUSED) {
					// agent can't connect to server, let's try again
					sleep := 1 * time.Second
					for i := uint(0); i < agent.retryCount; i++ {
						time.Sleep(sleep)
						err := agent.SendAllMetricsJSON()
						// continue if success
						if err == nil {
							break
						}
						// update sleep time
						sleep = time.Duration(sleep.Seconds() + float64(agent.retryStep))
					}
				}
				if curErrCount > agent.retryCount {
					return err
				}
			}
			curErrCount = 0
		}
	}
}


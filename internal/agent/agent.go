package agent

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/logger"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage"
	"net"
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

type Agent struct {
	addr           net.Addr
	metricsPoint   string
	metrics        map[string]metrics.MyMetrics
	store          storage.Storage
	pollInterval   time.Duration
	reportInterval time.Duration
	compression    bool
	logLevel       string
}

func (agent Agent) InitLogger(loggerLevel string) error {
	return logger.Initialize(loggerLevel)
}

func NewAgent(addr string, pollInterval, reportInterval int) (*Agent, error) {
	netAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	agent := &Agent{
		addr:           netAddr,
		metricsPoint:   "update",
		metrics:        make(map[string]metrics.MyMetrics),
		pollInterval:   time.Duration(pollInterval) * time.Second,
		reportInterval: time.Duration(reportInterval) * time.Second,
		logLevel:       "debug",
	}
	err = agent.InitLogger(agent.logLevel)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (agent Agent) Start() error {

	var (
		maxErrCount int
		curErrCount int
	)
	maxErrCount = 3
	curErrCount = 0
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
			agent.GetMetrics()
		case t2 := <-reportTicker.C:
			fmt.Println("- Sending:", t2.Format(time.TimeOnly))
			//agent.SendMetrics()
			err := agent.SendMetricsJSON()
			if err != nil {
				log.Debug().
					Msgf("%v error, while send metrics", err)
				curErrCount += 1
				if curErrCount > maxErrCount {
					return err
				}
			}
		}
	}
}

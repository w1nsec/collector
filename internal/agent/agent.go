package agent

import (
	"fmt"
	"github.com/w1nsec/collector/internal/metrics"
	"log"
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
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
	pollInterval   time.Duration
	reportInterval time.Duration
}

func NewAgent(addr string, pollInterval, reportInterval time.Duration) (*Agent, error) {
	netAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	agent := &Agent{
		addr:           netAddr,
		metricsPoint:   "update",
		metrics:        make(map[string]metrics.MyMetrics),
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
	}
	return agent, nil
}

func (agent Agent) SendMetrics() {
	if agent.metrics == nil {
		return
	}

	for mName, metric := range agent.metrics {
		url := fmt.Sprintf("http://%s/%s/%s/%s/%s", agent.addr.String(),
			agent.metricsPoint, metric.SendType, mName, metric.Value)
		fmt.Println(url)
		resp, err := http.Post(url, "text/plain", nil)
		// TODO handle error
		if err != nil {
			log.Println(err)
			continue
		}
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func (agent Agent) GetMetrics() {
	if agent.metrics == nil {
		agent.metrics = make(map[string]metrics.MyMetrics)
	}

	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	values := reflect.ValueOf(m)

	//metrics := make(map[string]interface{})
	for _, name := range usedMemStats {
		if values.FieldByName(name).IsValid() {
			if values.FieldByName(name).CanInt() {
				agent.metrics[name] = metrics.MyMetrics{
					Value:    strconv.FormatInt(values.FieldByName(name).Int(), 10),
					SendType: metrics.Counter,
				}
			}
			if values.FieldByName(name).CanUint() {
				agent.metrics[name] = metrics.MyMetrics{
					Value:    strconv.FormatUint(values.FieldByName(name).Uint(), 10),
					SendType: metrics.Counter,
				}
			}
			if values.FieldByName(name).CanFloat() {
				agent.metrics[name] = metrics.MyMetrics{
					Value:    strconv.FormatFloat(values.FieldByName(name).Float(), 'f', -1, 64),
					SendType: metrics.Gauge,
				}
			}
		}
	}

	// Addition metrics
	metric, ok := agent.metrics["PollCount"]
	if !ok {
		agent.metrics["PollCount"] = metrics.MyMetrics{
			Value:    "1",
			SendType: metrics.Counter,
		}
	} else {
		metric.AddVal(1)
		agent.metrics["PollCount"] = metric
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	randVal := float64(r.Intn(1000)) + r.Float64()
	agent.metrics["RandomValue"] = metrics.MyMetrics{
		Value:    strconv.FormatFloat(randVal, 'f', -1, 64),
		SendType: metrics.Gauge,
	}
}

func (agent Agent) Start() {
	// Receive and send for the first time
	fmt.Println("Receiving:", time.Now().Format(time.TimeOnly))
	agent.GetMetrics()
	fmt.Println("- Sending:", time.Now().Format(time.TimeOnly))
	agent.SendMetrics()

	pollTicker := time.NewTicker(agent.pollInterval)
	reportTicker := time.NewTicker(agent.reportInterval)
	for {
		select {
		case t1 := <-pollTicker.C:
			fmt.Println("Receiving:", t1.Format(time.TimeOnly))
			agent.GetMetrics()
		case t2 := <-reportTicker.C:
			fmt.Println("- Sending:", t2.Format(time.TimeOnly))
			agent.SendMetrics()
		}
	}
}

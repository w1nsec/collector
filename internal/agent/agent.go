package agent

import (
	"fmt"
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

type MyMetrics struct {
	value    string
	sendType string
}

func (m *MyMetrics) AddVal(n int) {
	// TODO parse error
	val, _ := strconv.ParseInt(m.value, 10, 64)
	val += int64(n)
	m.value = strconv.FormatInt(val, 10)
}

const (
	gauge   = "gauge"
	counter = "counter"
)

type Agent struct {
	addr    net.Addr
	metrics map[string]MyMetrics
}

func NewAgent(addr string) (*Agent, error) {
	netAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	agent := &Agent{
		addr:    netAddr,
		metrics: make(map[string]MyMetrics),
	}
	return agent, nil
}

func (agent Agent) SendMetrics() {
	if agent.metrics == nil {
		return
	}

	for mName, metric := range agent.metrics {
		url := fmt.Sprintf("http://%s/%s/%s/%s", agent.addr.String(), metric.sendType, mName, metric.value)
		fmt.Println(url)
		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			log.Println(err)
			continue
		}

		req.Header.Add("Content-Type", "text/plain")

	}
}

func (agent Agent) GetMetrics() {
	if agent.metrics == nil {
		agent.metrics = make(map[string]MyMetrics)
	}

	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	values := reflect.ValueOf(m)

	//metrics := make(map[string]interface{})
	for _, name := range usedMemStats {
		if values.FieldByName(name).IsValid() {
			if values.FieldByName(name).CanInt() {
				agent.metrics[name] = MyMetrics{
					value:    strconv.FormatInt(values.FieldByName(name).Int(), 10),
					sendType: counter,
				}
			}
			if values.FieldByName(name).CanUint() {
				agent.metrics[name] = MyMetrics{
					value:    strconv.FormatUint(values.FieldByName(name).Uint(), 10),
					sendType: counter,
				}
			}
			if values.FieldByName(name).CanFloat() {
				agent.metrics[name] = MyMetrics{
					value:    strconv.FormatFloat(values.FieldByName(name).Float(), 'f', -1, 64),
					sendType: gauge,
				}
			}
		}
	}

	// Addition metrics
	metric, ok := agent.metrics["PollCount"]
	if !ok {
		agent.metrics["PollCount"] = MyMetrics{
			value:    "1",
			sendType: counter,
		}
	} else {
		metric.AddVal(1)
		agent.metrics["PollCount"] = metric
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	randVal := float64(r.Intn(1000)) + r.Float64()
	agent.metrics["RandomValue"] = MyMetrics{
		value:    strconv.FormatFloat(randVal, 'f', -1, 64),
		sendType: gauge,
	}
}

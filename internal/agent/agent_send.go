package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	"io"
	"net/http"
	"strconv"
)

var (
	errConv = fmt.Errorf("error while converting \"mymetric\" to \"metric\"")
)

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
			//log.Println(err)
			continue
		}
		err = resp.Body.Close()
		if err != nil {
			//log.Println(err)
			continue
		}
	}
}

func (agent Agent) SendMetricsJSON() error {

	body, err := agent.generateRequest()

	if err != nil {
		return err
	}
	log.Info().RawJSON("body", body).
		Msg("Send: ")

	buffer := bytes.NewBuffer(body)
	//client := http.Client{
	//	Timeout: time.Duration(5),
	//}
	address := fmt.Sprintf("http://%s/%s/", agent.addr.String(), "update")
	//req, err := http.NewRequest("POST", address, buffer)
	//req, err := http.Post(address, "application/json", buffer)
	//if err != nil {
	//	return err
	//}

	//req.Header.Set("Content-type", "application/json")
	//resp, err := client.Do(req)
	//if err != nil {
	//	return err
	//}

	resp, err := http.Post(address, "application/json", buffer)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Info().
		RawJSON("body", body).
		Int("status", resp.StatusCode).
		Msg("Receive:")

	return nil
}

func convertOneMetric(name string, mymetric metrics.MyMetrics) (*metrics.Metrics, error) {
	metric := &metrics.Metrics{}
	switch mymetric.SendType {
	case metrics.Gauge:
		val, err := strconv.ParseFloat(mymetric.Value, 64)
		if err != nil {
			return nil, err
		}
		metric = metrics.NewGaugeMetric(name, metrics.Gauge, val)
		if metric == nil {
			return nil, errConv
		}
	case metrics.Counter:
		val, err := strconv.Atoi(mymetric.Value)
		if err != nil {
			return nil, err
		}
		metric = metrics.NewCounterMetric(name, metrics.Counter, int64(val))
		if metric == nil {
			return nil, errConv
		}
	}
	return metric, nil
}

func convertMetrics(agentMetric map[string]metrics.MyMetrics) ([]*metrics.Metrics, error) {
	newMetrics := make([]*metrics.Metrics, 0)
	for mName, mymetric := range agentMetric {
		metric, err := convertOneMetric(mName, mymetric)
		if err != nil {
			return nil, err
		}
		newMetrics = append(newMetrics, metric)
		//fmt.Println(metric)
	}
	return newMetrics, nil
}

func (agent Agent) generateRequest() ([]byte, error) {
	sendingMetrics, err := convertMetrics(agent.metrics)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(sendingMetrics)
	if err != nil {
		return nil, err
	}

	//buf := make([]byte, len(body))
	//buffer := bytes.NewBuffer(buf)
	//err = json.Indent(buffer, body, "", "  ")
	//if err != nil {
	//	return nil, err
	//}

	return body, err
}

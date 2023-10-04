package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	"io"
	"net/http"
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

func (agent Agent) SendOneMetricJSON(name string, mymetric metrics.MyMetrics) error {

	metric, err := metrics.ConvertMymetric2Metric(name, mymetric)
	if err != nil {
		return err
	}

	address := fmt.Sprintf("http://%s/%s/", agent.addr.String(), "update")

	body, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	log.Info().
		RawJSON("body", body).
		Msg("Send: ")

	buffer := bytes.NewBuffer(body)
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
		Msg("Receive: ")

	return nil
}

func (agent Agent) SendMetricsJSON() error {

	for name, mymetric := range agent.metrics {
		err := agent.SendOneMetricJSON(name, mymetric)
		if err != nil {
			log.Error().
				Err(err).Send()
		}
	}

	return nil
}

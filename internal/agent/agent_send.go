package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	locgzip "github.com/w1nsec/collector/internal/utils/compression/gzip"
	"github.com/w1nsec/collector/internal/utils/signing"
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

// SendOneMetricJSON send one metric in JSON format
func (agent Agent) SendOneMetricJSON(name string, mymetric metrics.MyMetrics) error {
	var (
		URL = "update"
	)

	metric, err := metrics.ConvertMymetric2Metric(name, mymetric)
	if err != nil {
		return err
	}

	//address := fmt.Sprintf("http://%s/%s/", agent.addr.String(), URL)

	body, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	log.Info().
		RawJSON("body", body).
		Msg("Send: ")

	headers := map[string]string{
		"content-type": "application/json",
	}
	return agent.SendData(body, headers, URL)
}

// SendMetricsJSON send collected metrics one by one in JSON format
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

// actual (iter14)
func (agent Agent) SendAllMetricsJSON() error {
	var (
		URL = "updates"
	)

	// encode metrics to json
	data := make([]byte, 0)
	buf := bytes.NewBuffer(data)
	encoder := json.NewEncoder(buf)
	all, err := agent.store.GetAllMetrics(context.TODO())
	if err != nil {
		return err
	}
	err = encoder.Encode(all)
	if err != nil {
		return err
	}

	// add http header
	headers := map[string]string{
		"content-type": "application/json",
	}

	return agent.SendData(buf.Bytes(), headers, URL)
}

func (agent Agent) SendData(data []byte, headers map[string]string, relURL string) error {
	// setup compression
	var buffer = bytes.NewBuffer(data)
	compressionStatus := false
	if agent.compression {
		compressed, err := locgzip.Compress(data)
		if err == nil {
			buffer = bytes.NewBuffer(compressed)
			compressionStatus = true
		}
	}

	// iter 14: add signing for body request
	agent.AddSigning(buffer.Bytes(), headers)

	// construct new request
	address := fmt.Sprintf("http://%s/%s/", agent.addr.String(), relURL)
	request, err := http.NewRequest(http.MethodPost, address, buffer)
	if err != nil {
		return err
	}

	// add headers for request
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// add compression header, if compression available
	if agent.compression && compressionStatus {
		request.Header.Set("content-encoding", "gzip")

		// decompress ??
		//request.Header.Set("accept-encoding", "gzip")
	}

	resp, err := agent.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Info().
		Str("url", relURL).
		RawJSON("body", body).
		Msg("Receive: ")

	return nil
}

func (agent Agent) AddSigning(data []byte, headers map[string]string) {
	if agent.secret == "" {
		return
	}
	sign := signing.CreateSigning(data, []byte(agent.secret))
	headers["HashSHA256"] = string(sign)

}

package agent

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config"
	"github.com/w1nsec/collector/internal/metrics"
	locgzip "github.com/w1nsec/collector/internal/utils/compression/gzip"
	"github.com/w1nsec/collector/internal/utils/signing"
	mycrypto "github.com/w1nsec/go-examples/crypto"
)

func (agent Agent) Send(job []*metrics.Metrics) error {
	switch agent.proto {
	case ProtoHTTP:
		return agent.SendBatch(job)
	case ProtoGRPC:
		return agent.SendRPCBatch(job)
	default:
		return fmt.Errorf("agent started with unsuported transport: %d", agent.proto)
	}
}

// SendData send data to server
// compress data before sending, add signing for body
func (agent Agent) SendData(data []byte, headers map[string]string, relURL string) error {
	// iter21: encrypt data
	if agent.pubKey != nil {
		var err error
		data, err = agent.Encrypt(data, headers)
		if err != nil {
			return err
		}
	}

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

	// iter24: add X-Real-IP header
	agent.AddRealIP(headers)

	// construct new request
	address := fmt.Sprintf("http://%s/%s/", agent.addr.String(), relURL)
	request, err := http.NewRequest(http.MethodPost, address, buffer)
	if err != nil {
		return fmt.Errorf("can't send request %v", err)
	}

	// add headers for request
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// add compression header, if compression available
	if agent.compression && compressionStatus {
		request.Header.Set("content-encoding", "gzip")
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
	if compressionStatus {
		log.Info().
			Str("url", relURL).
			Str("body", "gzip").
			Msg("Response:")
	} else {
		log.Info().
			Str("url", relURL).
			RawJSON("body", body).
			Msg("Response:")
	}

	return nil
}

// AddSigning func add signing, calculated for body to header "HashSHA256"
func (agent Agent) AddSigning(data []byte, headers map[string]string) {
	if agent.secret == "" {
		return
	}
	sign := signing.CreateSigning(data, []byte(agent.secret))
	headers[config.SignHeader] = string(sign)

}

// AddRealIP func add X-Real-IP header with agent ip addresses
func (agent Agent) AddRealIP(headers map[string]string) {
	if agent.realIP == nil {
		return
	}
	headers[config.RealIPHeader] = strings.Join(agent.realIP, ";")
}

// SendBatch prepare sending collected metrics to server
// encode slice of metrics to JSON format
// iter15 (such as SendAllMetricsJSON from iter14, but with args)
func (agent Agent) SendBatch(job []*metrics.Metrics) error {
	var (
		URL = "updates"
	)

	// encode metrics to json
	data := make([]byte, 0)
	buf := bytes.NewBuffer(data)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(job)
	if err != nil {
		return err
	}

	// add http header
	headers := map[string]string{
		"content-type": "application/json",
	}

	return agent.SendData(buf.Bytes(), headers, URL)
}

// Encrypt - this func encrypt body of message
// 1. generate AES key
// 2. encrypt AES key with RSA PublicKey
// 3. add encrypted AES key to headers value
// 4. encrypt message body with AES and return it
func (agent Agent) Encrypt(data []byte, headers map[string]string) (cData []byte, err error) {

	// generate random KEY for AES encryption
	secret, err := mycrypto.GenRandSlice(10)
	if err != nil {
		return nil, err
	}

	tmp := sha256.Sum256(secret)
	key := tmp[:]

	// encrypt AES key for sending with message
	keyEnc, err := rsa.EncryptPKCS1v15(rand.Reader, agent.pubKey, key)
	if err != nil {
		return nil, err
	}

	// add encrypted AES to header
	headers[config.CryptoHeader] = base64.StdEncoding.EncodeToString(keyEnc)

	// encrypt message body
	return mycrypto.EncryptAES(data, key)
}

// Legacy Code
// Deprecated
/*
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
*/

/*
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
*/

/*
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
*/

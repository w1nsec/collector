package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/w1nsec/collector/internal/metrics"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewJSONUpdateMetricsHandler(t *testing.T) {
	var usecase jsonUpdateUsecase
	got := NewJSONUpdateMetricsHandler(usecase)
	require.NotNil(t, got, "got nil value from constructor")
}

func (u *JSONusecase) UpdateMetrics(ctx context.Context, newMetrics []*metrics.Metrics) error {
	for _, m := range newMetrics {
		err := u.UpdateMetric(ctx, m)
		if err != nil {
			log.Error().Err(err).Send()
		}
	}
	return nil
}

func TestJSONUpdateMetricsHandler_ServeHTTP(t *testing.T) {

	gaugeKey := "validGauge"
	counterKey := "validCounter"
	gauge := float64(6.112)
	counter := int64(6112)
	mGauge := metrics.NewGaugeMetric(gaugeKey, gauge)
	mCounter := metrics.NewCounterMetric(counterKey, counter)

	existGaugeKey := "tmp_gauge"
	existCounterKey := "tmp_counter"
	existGaugeValue := float64(99.99)
	existCounterValue := int64(111)
	existGauge := metrics.NewGaugeMetric(existGaugeKey, existGaugeValue)
	existCounter := metrics.NewCounterMetric(existCounterKey, existCounterValue)

	existGaugeValueNew := float64(10.01)
	existCounterValueNew := int64(222)
	existGaugeNew := metrics.NewGaugeMetric(existGaugeKey, existGaugeValueNew)
	existCounterNew := metrics.NewCounterMetric(existCounterKey, existCounterValueNew)

	addNew1 := make([]*metrics.Metrics, 0)
	addNew1 = append(addNew1, mCounter)

	addNew2 := make([]*metrics.Metrics, 0)
	addNew2 = append(addNew2, mGauge)

	addNew3 := make([]*metrics.Metrics, 0)
	addNew3 = append(addNew3, mGauge, mCounter)

	addNew4 := make([]*metrics.Metrics, 0)
	addNew4 = append(addNew4, existGaugeNew, existCounterNew)

	reqCounter, err := json.Marshal(&addNew1)
	if err != nil {
		fmt.Printf("can't marshal to json: %v\n", err)
		return
	}

	reqGauge, err := json.Marshal(&addNew2)
	if err != nil {
		fmt.Printf("can't marshal to json: %v\n", err)
		return
	}

	reqAdd, err := json.Marshal(&addNew3)
	if err != nil {
		fmt.Printf("can't marshal to json: %v\n", err)
		return
	}

	reqUpd, err := json.Marshal(&addNew4)
	if err != nil {
		fmt.Printf("can't marshal to json: %v\n", err)
		return
	}

	fmt.Println(string(reqCounter))
	fmt.Println(string(reqGauge))
	fmt.Println(string(reqAdd))
	fmt.Println(string(reqUpd))

	type args struct {
		method      string
		contentType map[string]string
		requestBody []byte
		respBody    string
		resStatus   int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test GET",
			args: args{
				method: http.MethodGet,
				contentType: map[string]string{
					"content-type": "application/json",
				},
				requestBody: nil,
				resStatus:   http.StatusMethodNotAllowed,
			},
		},
		{
			name: "Test empty POST without Content-Type",
			args: args{
				method:      http.MethodPost,
				requestBody: nil,
				resStatus:   http.StatusInternalServerError,
			},
		},
		{
			name: "Test empty POST",
			args: args{
				method: http.MethodPost,
				contentType: map[string]string{
					"content-type": "application/json",
				},
				requestBody: nil,
				resStatus:   http.StatusInternalServerError,
			},
		},

		{
			name: "Test POST, add new Counter",
			args: args{
				method: http.MethodPost,
				contentType: map[string]string{
					"content-type": "application/json",
				},
				requestBody: reqCounter,
				respBody:    `[{"id":"validCounter","type":"counter","delta":6112}]`,
				resStatus:   http.StatusOK,
			},
		},
		{
			name: "Test POST, add new Gauge",
			args: args{
				method: http.MethodPost,
				contentType: map[string]string{
					"content-type": "application/json",
				},
				requestBody: reqGauge,
				respBody:    `[{"id":"validGauge","type":"gauge","value":6.112}]`,
				resStatus:   http.StatusOK,
			},
		},
		{
			name: "Test POST, add new Counter + Gauge",
			args: args{
				method: http.MethodPost,
				contentType: map[string]string{
					"content-type": "application/json",
				},
				requestBody: reqAdd,
				respBody:    `[{"id":"validGauge","type":"gauge","value":6.112},{"id":"validCounter","type":"counter","delta":6112}]`,
				resStatus:   http.StatusOK,
			},
		},
		{
			name: "Test POST, update Counter + Gauge",
			args: args{
				method: http.MethodPost,
				contentType: map[string]string{
					"content-type": "application/json",
				},
				requestBody: reqUpd,
				respBody:    `[{"id":"tmp_gauge","type":"gauge","value":10.01},{"id":"tmp_counter","type":"counter","delta":333}]`,
				resStatus:   http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var usecase JSONusecase
			usecase.m = make(map[string]*metrics.Metrics)
			usecase.m[existGaugeKey] = existGauge
			usecase.m[existCounterKey] = existCounter

			h := &JSONUpdateMetricsHandler{
				usecase: &usecase,
			}

			fmt.Println(existGaugeKey, ":   ", usecase.m[existGaugeKey])
			fmt.Println(existCounterKey, ":   ", usecase.m[existCounterKey])

			buf := bytes.NewBuffer(tt.args.requestBody)
			req := httptest.NewRequest(tt.args.method, "localhost:8000/updates/", buf)
			if tt.args.contentType != nil {
				for k, v := range tt.args.contentType {
					req.Header.Add(k, v)
				}
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, res.StatusCode, tt.args.resStatus, "status code not valid")
			if res.StatusCode == http.StatusOK {
				body, err := io.ReadAll(res.Body)
				require.NoError(t, err, "can't read body response")

				sBody := strings.TrimSpace(string(body))
				require.Equal(t, sBody, strings.TrimSpace(tt.args.respBody))
			}

		})
	}
}

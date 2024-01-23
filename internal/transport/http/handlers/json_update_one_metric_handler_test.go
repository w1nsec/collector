package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/w1nsec/collector/internal/metrics"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (u *JSONusecase) UpdateMetric(ctx context.Context, newMetric *metrics.Metrics) error {
	v, ok := u.m[newMetric.ID]
	// metric not exist, add it
	if !ok {
		u.m[newMetric.ID] = newMetric
		return nil
	}

	// metric already exist, update it
	switch v.MType {
	case metrics.Gauge:
		*u.m[newMetric.ID].Value = *newMetric.Value
		return nil
	case metrics.Counter:
		*u.m[newMetric.ID].Delta += *newMetric.Delta
		return nil
	}
	return fmt.Errorf("metric type not supported")
}

func TestJSONUpdateOneMetricHandler_ServeHTTP(t *testing.T) {
	usecase := &JSONusecase{
		m: make(map[string]*metrics.Metrics),
	}

	addr := "localhost:8000/update/"
	ct := "application/json"

	gaugeKey := "validGauge"
	counterKey := "validCounter"
	gauge := float64(6.112)
	counter := int64(6112)
	mGauge := metrics.NewGaugeMetric(gaugeKey, gauge)
	mCounter := metrics.NewCounterMetric(counterKey, counter)

	usecase.m[gaugeKey] = mGauge
	usecase.m[counterKey] = mCounter

	existCounter := metrics.NewCounterMetric(counterKey, 111)
	exist1, err := json.Marshal(&existCounter)
	if err != nil {
		fmt.Printf("can't marshal to json: %v\n", err)
		return
	}
	existGauge := metrics.NewGaugeMetric(gaugeKey, 999.99)
	exist2, err := json.Marshal(&existGauge)
	if err != nil {
		fmt.Printf("can't marshal to json: %v\n", err)
		return
	}

	counterKey1 := "newCounter"
	gaugeKey1 := "newGauge"
	newCounter := metrics.NewCounterMetric(counterKey1, 111)
	new1, err := json.Marshal(&newCounter)
	if err != nil {
		fmt.Printf("can't marshal to json: %v\n", err)
		return
	}
	newGauge := metrics.NewGaugeMetric(gaugeKey1, 999.99)
	new2, err := json.Marshal(&newGauge)
	if err != nil {
		fmt.Printf("can't marshal to json: %v\n", err)
		return
	}
	type args struct {
		method      string
		contentType string
		body        []byte
		status      int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test update exist counter",
			args: args{
				method:      http.MethodPost,
				body:        exist1,
				contentType: ct,
				status:      http.StatusOK,
			},
		},
		{
			name: "Test update exist gauge",
			args: args{
				method:      http.MethodPost,
				body:        exist2,
				contentType: ct,
				status:      http.StatusOK,
			},
		},
		{
			name: "Test update new counter",
			args: args{
				method:      http.MethodPost,
				body:        new1,
				contentType: ct,
				status:      http.StatusOK,
			},
		},
		{
			name: "Test update new gauge",
			args: args{
				method:      http.MethodPost,
				body:        new2,
				contentType: ct,
				status:      http.StatusOK,
			},
		},
		{
			name: "Test invalid body (empty)",
			args: args{
				method: http.MethodPost,
				body:   []byte(""),
				//body:   []byte(`{"id":"validCounter","type":"counter","delta":111}`),
				contentType: ct,
				status:      http.StatusInternalServerError,
			},
			wantErr: true,
		},
		{
			name: "Test invalid json in body",
			args: args{
				method:      http.MethodPost,
				body:        []byte(`{"id":"validCounter","type":"counter","delta:111}`),
				contentType: ct,
				status:      http.StatusInternalServerError,
			},
			wantErr: true,
		},
		{
			name: "Test invalid json in body (no metric value)",
			args: args{
				method:      http.MethodPost,
				body:        []byte(`{"id":"validCounter","type":"counter"}`),
				contentType: ct,
				status:      http.StatusInternalServerError,
			},
			wantErr: true,
		},

		{
			name: "Test invalid content-type",
			args: args{
				method:      http.MethodPost,
				body:        []byte(`{"id":"validCounter","type":"counter","delta":111}`),
				contentType: "text/plain",
				status:      http.StatusInternalServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &JSONUpdateOneMetricHandler{
				usecase: usecase,
			}
			m := &metrics.Metrics{}
			err := json.Unmarshal(tt.args.body, m)
			if err != nil && (err != nil) != tt.wantErr {
				t.Errorf("can't unmarshal body to metric struct: %v", err)
				return
			}
			oldM, err := h.usecase.GetMetric(context.Background(), m.ID, m.MType)
			var oldD int64
			if err != nil {
				//t.Errorf("can't get old metric value: %v", err)
				oldD = 0
			} else {
				if oldM.Delta != nil {
					oldD = *oldM.Delta
				}
			}

			body := bytes.NewBuffer(tt.args.body)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tt.args.method, addr, body)
			if tt.args.contentType != "" {
				req.Header.Add("content-type", tt.args.contentType)
			}
			h.ServeHTTP(rec, req)
			res := rec.Result()
			defer res.Body.Close()

			require.Equal(t, tt.args.status, res.StatusCode)

			if res.StatusCode == http.StatusOK {
				gotM, err := h.usecase.GetMetric(context.Background(), m.ID, m.MType)
				if err != nil {
					t.Errorf("can't get new metric value: %v", err)
					return
				}
				require.Equal(t, m.ID, gotM.ID)
				require.Equal(t, m.MType, gotM.MType)
				switch gotM.MType {
				case metrics.Counter:
					//fmt.Printf("old: %d | new: %d | got: %d", *oldM.Delta, *m.Delta, *gotM.Delta)
					expect := *m.Delta + oldD
					got := *gotM.Delta
					require.Equal(t, expect, got)
				case metrics.Gauge:
					require.Equal(t, *m.Value, *gotM.Value)
				default:
					t.Errorf("wrong metric type: %s", m.MType)
					return
				}

			}
		})
	}
}

func TestNewJSONUpdateOneMetricHandler(t *testing.T) {
	var usecase JSONusecase
	got := NewJSONUpdateOneMetricHandler(&usecase)
	require.NotNil(t, got, "got nil value from constructor")
}

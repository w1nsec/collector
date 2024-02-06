package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage/memstorage"
)

func TestGetMetricHandler_ServeHTTP(t *testing.T) {
	store := memstorage.NewMemStorage()
	gauge := float64(6.112)
	counter := int64(6112)
	mGauge := metrics.NewGaugeMetric("validGauge", gauge)
	mCounter := metrics.NewCounterMetric("validCounter", counter)

	err := store.AddMetric(context.Background(), mGauge)
	if err != nil {
		fmt.Println("can't add counter metric to storage")
		return
	}
	err = store.AddMetric(context.Background(), mCounter)
	if err != nil {
		fmt.Println("can't add gauge metric to storage")
		return
	}

	type args struct {
		method string
		key    string
		value  string
		status int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test Valid Counter metric",
			args: args{
				method: http.MethodGet,
				key:    metrics.Counter,
				value:  "validCounter",
				status: http.StatusOK,
			},
		},
		{
			name: "Test Valid Gauge metric",
			args: args{
				method: http.MethodGet,
				key:    metrics.Gauge,
				value:  "validGauge",
				status: http.StatusOK,
			},
		},
		{
			name: "Test Invalid Counter metric",
			args: args{
				method: http.MethodGet,
				key:    metrics.Counter,
				value:  "invalidCounter",
				status: http.StatusNotFound,
			},
		},
		{
			name: "Test Invalid Gauge metric",
			args: args{
				method: http.MethodGet,
				key:    metrics.Gauge,
				value:  "invalidGauge",
				status: http.StatusNotFound,
			},
		},
		{
			name: "Test Invalid method",
			args: args{
				method: http.MethodPost,
				key:    metrics.Gauge,
				value:  "validGauge",
				status: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &GetMetricHandler{
				usecase: store,
			}
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(tt.args.method, "http://localhost:8000/{mType}/{mName}", nil)

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("mType", tt.args.key)
			ctx.URLParams.Add("mName", tt.args.value)
			req = req.WithContext(context.WithValue(req.Context(),
				chi.RouteCtxKey, ctx))

			h.ServeHTTP(recorder, req)
			res := recorder.Result()
			defer res.Body.Close()

			require.Equal(t, tt.args.status, res.StatusCode)

			if res.StatusCode == http.StatusOK {
				if ct := res.Header.Get("content-type"); ct != "text/plain" {
					t.Errorf("response has wrong content-type: %s, want: \"text/plain\"", ct)
					return
				}
			}
		})
	}
}

type getAll struct{}

func (g getAll) String(ctx context.Context) string { return "OK" }

func TestGetMetricsHandler_ServeHTTP(t *testing.T) {
	usecase := &getAll{}
	type args struct {
		method string
		status int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test POST request",
			args: args{
				method: http.MethodPost,
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Test GET request",
			args: args{
				method: http.MethodGet,
				status: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &GetMetricsHandler{
				usecase: usecase,
			}

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(tt.args.method, "http://localhost:8000/", nil)
			h.ServeHTTP(recorder, req)

			res := recorder.Result()
			defer res.Body.Close()

			require.Equal(t, tt.args.status, res.StatusCode)

			if tt.args.method == http.MethodGet {
				if ct := res.Header.Get("content-type"); ct != "text/html" {
					t.Errorf("wrong conten-type: %s, want: \"text/html\"", ct)
					return
				}
			}

		})
	}
}

type mockGetMetric struct {
	m map[string]*metrics.Metrics
	//jsonGetMetricUsecase
}

func (m mockGetMetric) GetMetricString(ctx context.Context, mType string, mName string) string {
	return "all is good"
}

func TestNewGetMetricHandler(t *testing.T) {
	var usecase mockGetMetric
	gotHdl := NewGetMetricHandler(usecase)
	require.NotNil(t, gotHdl, "got nil value from constructor")
}

func TestNewGetMetricsHandler(t *testing.T) {
	var usecase getAllMetricsUsecase
	gotHdl := NewGetMetricsHandler(usecase)
	require.NotNil(t, gotHdl, "got nil value from constructor")
}

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage/memstorage"
)

type JSONusecase struct {
	m map[string]*metrics.Metrics
	//jsonGetMetricUsecase
}

func (u JSONusecase) GetMetric(ctx context.Context, mName string, mType string) (*metrics.Metrics, error) {
	if u.m == nil {
		u.m = make(map[string]*metrics.Metrics)
	}
	m1 := metrics.NewCounterMetric("test1", 111)
	m2 := metrics.NewGaugeMetric("test2", 333.1)

	u.m[m1.ID] = m1
	u.m[m2.ID] = m2

	m, ok := u.m[mName]
	if !ok {
		return nil, fmt.Errorf("metric with name:%s not found", mName)
	}

	if m.MType != mType {
		return nil, fmt.Errorf("metric with name:%s exist, but wrong type, got: %s, want: %s", mName, mType, m.MType)
	}
	return m, nil
}

func TestJSONGetMetricHandler_ServeHTTP(t *testing.T) {
	usecase := &JSONusecase{}
	addr := "http://localhost:8000/value/"
	m1 := metrics.NewCounterMetric("test1", 111)
	m1Err := metrics.NewCounterMetric("test13", 1111)
	m2 := metrics.NewGaugeMetric("test2", 333.1)
	m2Err := metrics.NewGaugeMetric("test14", 333.1111)

	/*
		body1, err := json.Marshal(&m1)
		if err != nil {
			fmt.Printf("can't marshal metric to json: %v\n", err)
			return
		}
		body2, err := json.Marshal(&m2)
		if err != nil {
			fmt.Printf("can't marshal metric to json: %v\n", err)
			return
		}
	*/
	type args struct {
		method string
		metric *metrics.Metrics
		//body   []byte
		headers map[string]string
		status  int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test get valid counter",
			args: args{
				method: http.MethodPost,
				metric: m1,
				headers: map[string]string{
					"content-type": "application/json",
				},
				status: http.StatusOK,
			},
		},
		{
			name: "Test get not exist counter",
			args: args{
				method: http.MethodPost,
				metric: m1Err,
				headers: map[string]string{
					"content-type": "application/json",
				},
				status: http.StatusNotFound,
			},
		},
		{
			name: "Test get valid gauge",
			args: args{
				method: http.MethodPost,
				metric: m2,
				headers: map[string]string{
					"content-type": "application/json",
				},
				status: http.StatusOK,
			},
		},
		{
			name: "Test get not exist gauge",
			args: args{
				method: http.MethodPost,
				metric: m2Err,
				headers: map[string]string{
					"content-type": "application/json",
				},
				status: http.StatusNotFound,
			},
		},
		{
			name: "Test GET method",
			args: args{
				method: http.MethodGet,
				headers: map[string]string{
					"content-type": "application/json",
				},
				status: http.StatusInternalServerError,
			},
		},
		{
			name: "Test not application/json content-type",
			args: args{
				method: http.MethodGet,
				headers: map[string]string{
					"content-type": "wrong",
				},
				status: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewJSONGetMetricHandler(usecase)

			data, err := json.Marshal(tt.args.metric)
			if err != nil {
				t.Errorf("can't marshal metric to json: %v\n", err)
				return
			}
			body := bytes.NewBuffer(data)
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(tt.args.method, addr, body)
			for k, v := range tt.args.headers {
				req.Header.Add(k, v)
			}
			h.ServeHTTP(recorder, req)
			res := recorder.Result()
			defer res.Body.Close()

			require.Equal(t, tt.args.status, res.StatusCode)

			if res.StatusCode != http.StatusOK {
				return
			}

			if ct := res.Header.Get("content-type"); ct != "application/json" {
				t.Errorf("wrong content-type: %s, want: application/json\n", ct)
				return
			}

			var m *metrics.Metrics
			err = json.NewDecoder(res.Body).Decode(&m)
			if err != nil {
				t.Errorf("can't unmarshal metric to json: %v\n", err)
				return
			}

			require.Equal(t, tt.args.metric.ID, m.ID)
			require.Equal(t, tt.args.metric.MType, m.MType)
			if m.Delta != nil && tt.args.metric.Delta != nil {
				require.Equal(t, *tt.args.metric.Delta, *m.Delta)
			}
			if m.Value != nil && tt.args.metric.Value != nil {
				require.Equal(t, *tt.args.metric.Value, *m.Value)
			}

		})
	}
}

func ExampleJSONGetMetricHandler_ServeHTTP() {
	addr := "localhost:8000"
	path := "/value/"
	store := memstorage.NewMemStorage()
	m1 := metrics.NewCounterMetric("test1", 111)
	m2 := metrics.NewGaugeMetric("test2", 333.1)
	err := store.AddMetric(context.Background(), m1)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(context.Background(), m2)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}

	updateCounter := NewJSONGetMetricHandler(store)
	router := chi.NewRouter()
	router.Post(path, updateCounter.ServeHTTP)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go http.ListenAndServe(addr, router)
	time.Sleep(time.Second)
	bodyM := &metrics.Metrics{
		ID:    m1.ID,
		MType: m1.MType,
	}

	data, err := json.Marshal(bodyM)
	if err != nil {
		fmt.Printf("can't marshal body to json: %v", err)
		return
	}
	body := bytes.NewBuffer(data)
	url := fmt.Sprintf("http://%s%s", addr, path)
	client := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		fmt.Printf("can't create request: %v\n", err)
		return
	}
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("can't send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response SC:", resp.StatusCode)

	var respMetric *metrics.Metrics
	err = json.NewDecoder(resp.Body).Decode(&respMetric)
	if err != nil {
		fmt.Printf("can't decode body from json: %v", err)
		return
	}

	fmt.Println("Metric:\n", respMetric)
	wg.Done()
}

func TestExampleJSONGetMetricHandler_ServeHTTP(t *testing.T) {
	ExampleJSONGetMetricHandler_ServeHTTP()
}

func TestNewJSONGetMetricHandler(t *testing.T) {
	var usecase JSONusecase
	got := NewJSONGetMetricHandler(usecase)
	require.NotNil(t, got, "got nil value from constructor")
}

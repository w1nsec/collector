package handlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage/memstorage"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestUpdateCountersHandler_ServeHTTP(t *testing.T) {
	store := memstorage.NewMemStorage()

	type args struct {
		method string
		body   string
		name   string
		value  string
	}
	tests := []struct {
		name    string
		usecase updateCountersUsecase
		args    args
		status  int
	}{
		// TODO: Add test cases.
		{
			name:    "Test Valid value",
			usecase: store,
			args: args{
				method: http.MethodPost,
				body:   "",
				name:   "testname",
				value:  "123",
			},
			status: http.StatusOK,
		},
		{
			name:    "Test Invalid value",
			usecase: store,
			args: args{
				method: http.MethodPost,
				body:   "",
				name:   "testname",
				value:  "wrong",
			},
			status: http.StatusBadRequest,
		},
		{
			name:    "Test Valid value with body",
			usecase: store,
			args: args{
				method: http.MethodPost,
				body:   "this is simple request body",
				name:   "testname",
				value:  "123",
			},
			status: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UpdateCountersHandler{
				counterUsecase: store,
			}
			body := bytes.NewBufferString(tt.args.body)
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(tt.args.method,
				fmt.Sprintf("http://localhost:8000/counter/%s/%s", tt.args.name, tt.args.value),
				body)

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("name", tt.args.name)
			ctx.URLParams.Add("value", tt.args.value)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			h.ServeHTTP(recorder, req)
			result := recorder.Result()
			defer result.Body.Close()

			require.Equal(t, tt.status, result.StatusCode, fmt.Errorf("wrong statuscod got=%d, want=%d",
				result.StatusCode, tt.status))

		})
	}
}

func ExampleUpdateCountersHandler_ServeHTTP() {
	// metric params
	mName := "testmetric"
	mValue := "123"

	addr := "localhost:8000"
	path := "/update/counter/{name}/{value}"
	store := memstorage.NewMemStorage()

	updateCounter := NewUpdateCountersHandler(store)
	router := chi.NewRouter()
	router.Post(path, updateCounter.ServeHTTP)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go http.ListenAndServe(addr, router)

	url := fmt.Sprintf("http://%s/update/counter/%s/%s", addr, mName, mValue)
	client := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Printf("can't create request: %v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("can't send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response SC:", resp.StatusCode)

	m, err := store.GetMetric(context.Background(), mName, metrics.Counter)
	if err != nil {
		fmt.Printf("Metric not found")
		return
	}

	fmt.Println("Metric:", m)
	wg.Done()
}

func TestExampleUpdateCountersHandler_ServeHTTP(t *testing.T) {
	ExampleUpdateCountersHandler_ServeHTTP()
}

type mockUpdateCounters struct {
}

func (m mockUpdateCounters) UpdateCounters(ctx context.Context, name string, value int64) error {
	return nil
}

func TestNewUpdateCountersHandler(t *testing.T) {
	var usecase mockUpdateCounters
	got := NewUpdateCountersHandler(usecase)
	require.NotNil(t, got, "got nil value from constructor")
}

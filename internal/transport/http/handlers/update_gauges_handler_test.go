package handlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/w1nsec/collector/internal/storage/memstorage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateGaugeHandler_ServeHTTP(t *testing.T) {
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
			name:    "Test Valid int value",
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
			name:    "Test Valid float value",
			usecase: store,
			args: args{
				method: http.MethodPost,
				body:   "",
				name:   "testname",
				value:  "123.11",
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
			h := &UpdateGaugeHandler{
				gaugeUsecase: store,
			}
			body := bytes.NewBufferString(tt.args.body)
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(tt.args.method,
				fmt.Sprintf("http://localhost:8000/gauge/%s/%s", tt.args.name, tt.args.value),
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

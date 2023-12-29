package handlers

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBadRequest(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		statuscode int
	}{
		{
			name:       "Test Bad Request GET",
			method:     http.MethodGet,
			statuscode: http.StatusBadRequest,
		},
		{
			name:       "Test Bad Request POST",
			method:     http.MethodPost,
			statuscode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respWr := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, "http://localhost:8000/wrongrequest", nil)
			BadRequest(respWr, req)

			res := respWr.Result()
			defer res.Body.Close()

			require.Equal(t, tt.statuscode, res.StatusCode)
		})
	}
}

func TestNotFoundHandle(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		statuscode int
	}{
		{
			name:       "Test NotFound GET",
			method:     http.MethodGet,
			statuscode: http.StatusNotFound,
		},
		{
			name:       "Test NotFound POST",
			method:     http.MethodPost,
			statuscode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respWr := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, "http://localhost:8000/unknownPath", nil)
			NotFoundHandle(respWr, req)

			res := respWr.Result()
			defer res.Body.Close()

			require.Equal(t, tt.statuscode, res.StatusCode)
		})
	}
}

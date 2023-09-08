package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/w1nsec/collector/internal/memstorage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricsHandle(t *testing.T) {
	updateURL := "/update"
	tests := []struct {
		name       string
		store      memstorage.Storage
		haveUrl    string
		statusCode int
		//want http.HandlerFunc

	}{
		// TODO: Add test cases.
		{
			name:       "Wrong url 1",
			store:      memstorage.NewMemStorage(),
			haveUrl:    "/aaaaa",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Wrong url 2",
			store:      memstorage.NewMemStorage(),
			haveUrl:    updateURL + "/aaaaa",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Wrong url 3",
			store:      memstorage.NewMemStorage(),
			haveUrl:    updateURL + "/counter/sys",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Wrong Type",
			store:      memstorage.NewMemStorage(),
			haveUrl:    updateURL + "/aaaaa/sys/123",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Wrong Type Counter",
			store:      memstorage.NewMemStorage(),
			haveUrl:    updateURL + "/counter/sys/213.214",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Wrong Type Gauge",
			store:      memstorage.NewMemStorage(),
			haveUrl:    updateURL + "/wrongtype/sys/213",
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.haveUrl, nil)
			w := httptest.NewRecorder()
			hf := UpdateMetricsHandle(tt.store)
			hf(w, r)
			resp := w.Result()
			assert.Equal(t, resp.StatusCode, tt.statusCode)

			// TODO write test for storage
			//assert.Equal(t, resp.)
		})
	}
}

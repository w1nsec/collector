package handlers

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type wrongStore struct{}

func (s wrongStore) CheckStorage() error {
	return fmt.Errorf("wrong store")
}

type validStore struct{}

func (s validStore) CheckStorage() error {
	return nil
}

func TestCheckDBConnectionHandler_ServeHTTP(t *testing.T) {
	type args struct {
		method  string
		storage checkStorageUsecase
		status  int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "Test valid storage",
			args: args{
				method:  http.MethodGet,
				storage: &validStore{},
				status:  http.StatusOK,
			},
		}, {
			name: "Test valid storage POST",
			args: args{
				method:  http.MethodPost,
				storage: &validStore{},
				status:  http.StatusOK,
			},
		},
		{
			name: "Test nil storage",
			args: args{
				method:  http.MethodGet,
				storage: &wrongStore{},
				status:  http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CheckDBConnectionHandler{
				checkUsecase: tt.args.storage,
			}

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(tt.args.method, "http://localhost:8000/ping", nil)
			require.NoError(t, err)

			h.ServeHTTP(recorder, req)
			res := recorder.Result()
			defer res.Body.Close()

			require.Equal(t, tt.args.status, res.StatusCode)

		})
	}
}

func TestNewCheckDBConnectionHandler(t *testing.T) {
	var usecase checkStorageUsecase
	gotHdl := NewCheckDBConnectionHandler(usecase)
	require.NotNil(t, gotHdl, "got nil value from constructor")
}

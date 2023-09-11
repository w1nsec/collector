package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/w1nsec/collector/internal/memstorage"
	"net/http"
	"testing"
)

func TestNewMetricServer(t *testing.T) {
	type args struct {
		addr  string
		store memstorage.Storage
		mux   *http.ServeMux
	}
	defaultArgs := args{
		addr:  "127.0.0.1:8080",
		store: memstorage.NewMemStorage(),
		mux:   http.NewServeMux(),
	}

	tests := []struct {
		name    string
		args    args
		want    *MetricServer
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "resolve error",
			args:    args{"1271.0.0.1:8080", nil, nil},
			wantErr: true,
		},
		{
			name:    "Seems good",
			args:    defaultArgs,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetricServerWithParams(tt.args.addr, tt.args.store, tt.args.mux)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetricServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			assert.Equal(t, got.Addr, tt.args.addr)
			assert.Equal(t, got.Store, tt.args.store)
			assert.Equal(t, got.Handler, tt.args.mux)

		})
	}
}

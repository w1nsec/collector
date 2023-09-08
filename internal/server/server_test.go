package server

import (
	"github.com/w1nsec/collector/internal/memstorage"
	"net/http"
	"reflect"
	"testing"
)

func TestNewMetricServer(t *testing.T) {
	type args struct {
		addr  string
		store memstorage.Storage
		mux   *http.ServeMux
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
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetricServer(tt.args.addr, tt.args.store, tt.args.mux)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetricServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetricServer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

package agent

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/w1nsec/collector/internal/metrics"
	"reflect"
	"testing"
)

func Test_convertOneMetric(t *testing.T) {
	type args struct {
		name     string
		mymetric metrics.MyMetrics
	}

	var (
		test1 = int64(10)
		test2 = float64(10.1)
	)

	tests := []struct {
		name    string
		args    args
		want    *metrics.Metrics
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test Counter",
			args: args{
				name: "Test-Counter-metric",
				mymetric: metrics.MyMetrics{
					Value:    "10",
					SendType: metrics.Counter,
				},
			},
			want: &metrics.Metrics{
				ID:    "Test-Counter-metric",
				MType: "counter",
				Delta: &test1,
			},
			wantErr: false,
		},
		{
			name: "test Gauge",
			args: args{
				name: "Test-Gauge-metric",
				mymetric: metrics.MyMetrics{
					Value:    "10.1",
					SendType: metrics.Gauge,
				},
			},
			want: &metrics.Metrics{
				ID:    "Test-Gauge-metric",
				MType: "gauge",
				Value: &test2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertOneMetric(tt.args.name, tt.args.mymetric)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertOneMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
			// TODO
			if tt.args.name != got.ID {
				t.Errorf("convertOneMetric() wrong metric name, got = %v, want %v", got.ID, tt.args.name)
				return
			}
			if tt.args.mymetric.SendType != got.MType {
				t.Errorf("convertOneMetric() wrong metric type, got = %v, want %v", got.MType, tt.args.mymetric.SendType)
				return
			}
			val := reflect.ValueOf(tt.args.mymetric.Value)
			if val.CanFloat() && got.MType == "gauge" {
				require.Equal(t, *got.Value, val.Float())
			}
			if val.CanInt() && got.MType == "counter" {
				require.Equal(t, *got.Value, val.Int())
			}

		})
	}
}

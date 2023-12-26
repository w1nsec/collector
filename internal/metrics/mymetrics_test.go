package metrics

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ConvertMymetric2Metric(t *testing.T) {
	type args struct {
		name     string
		mymetric MyMetrics
	}

	var (
		test1 = int64(10)
		test2 = float64(10.1)
	)

	tests := []struct {
		name    string
		args    args
		want    *Metrics
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test Counter",
			args: args{
				name: "Test-Counter-metric",
				mymetric: MyMetrics{
					Value:    "10",
					SendType: Counter,
				},
			},
			want: &Metrics{
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
				mymetric: MyMetrics{
					Value:    "10.1",
					SendType: Gauge,
				},
			},
			want: &Metrics{
				ID:    "Test-Gauge-metric",
				MType: "gauge",
				Value: &test2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertMymetric2Metric(tt.args.name, tt.args.mymetric)
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

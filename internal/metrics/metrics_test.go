package metrics

import (
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestNewGaugeMetric(t *testing.T) {
	type args struct {
		name  string
		mType string
		value float64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "Test 1",
			args: args{
				name:  "newname",
				mType: Gauge,
				value: 1123.11,
			},
		},
		{
			name: "Test 2 max float64 value",
			args: args{
				name:  "newname",
				mType: Gauge,
				value: math.MaxFloat64,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGaugeMetric(tt.args.name, tt.args.value)
			require.Equal(t, tt.args.mType, got.MType)
			require.Equal(t, tt.args.name, got.ID)
			require.Equal(t, tt.args.value, *got.Value)
		})
	}
}

func TestNewCounterMetric(t *testing.T) {
	type args struct {
		name  string
		mType string
		value int64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "Test 1",
			args: args{
				name:  "newname",
				mType: Counter,
				value: 123,
			},
		},
		{
			name: "Test 2 max int64 value",
			args: args{
				name:  "newname",
				mType: Counter,
				value: math.MaxInt64,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCounterMetric(tt.args.name, tt.args.value)
			require.Equal(t, tt.args.mType, got.MType)
			require.Equal(t, tt.args.name, got.ID)
			require.Equal(t, tt.args.value, *got.Delta)
		})
	}
}

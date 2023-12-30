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

func TestDelete(t *testing.T) {
	m1 := NewCounterMetric("1111", 1)
	m2 := NewCounterMetric("2222", 2)
	m3 := NewCounterMetric("3333", 3)
	m4 := NewCounterMetric("4444", 4)

	sl1 := make([]*Metrics, 4)
	sl1[0] = m1
	sl1[1] = m2
	sl1[2] = m3
	sl1[3] = m4

	sl2 := make([]*Metrics, 4)
	copy(sl2, sl1)
	type args struct {
		metrics []*Metrics
		ind     int
	}
	tests := []struct {
		name  string
		args  args
		found bool
	}{
		{
			name: "Test index > len",
			args: args{
				metrics: sl1,
				ind:     len(sl1) + 2,
			},
			found: true,
		},
		{
			name: "Test index < len",
			args: args{
				metrics: sl1,
				ind:     len(sl1) - 2,
			},
			found: false,
		},
		{
			name: "Test index = last elem",
			args: args{
				metrics: sl2,
				ind:     len(sl2) - 1,
			},
			found: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			var deletingItem = new(Metrics)

			if len(tt.args.metrics) > tt.args.ind {
				if tt.found == false {

					deletingItem.ID = tt.args.metrics[tt.args.ind].ID
					deletingItem.MType = tt.args.metrics[tt.args.ind].MType
					d := *tt.args.metrics[tt.args.ind].Delta
					deletingItem.Delta = &d
				}
				newSl := Delete(tt.args.metrics, tt.args.ind)
				if tt.found == false {
					for _, m := range newSl {
						if m.ID == deletingItem.ID &&
							*m.Delta == *deletingItem.Delta {
							found = true
							break
						}
					}
				} else {
					found = true
				}

				require.Equal(t, tt.found, found)
				if found == false {
					require.Equal(t, len(tt.args.metrics)-1, len(newSl))
				}
			}
		})
	}
}

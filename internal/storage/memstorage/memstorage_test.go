package memstorage

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/w1nsec/collector/internal/metrics"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestMemStorage_Init(t *testing.T) {
	tests := []struct {
		name    string
		storage *MemStorage
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test checking maps for nil",
			storage: &MemStorage{
				mutex: &sync.RWMutex{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.storage.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() got error = %v, but wantErr %v", err, tt.wantErr)
			}

			require.NotNil(t, tt.storage.dataGauges)
			require.NotNil(t, tt.storage.dataCounters)
		})
	}
}

func TestMemStorage_CheckStorage(t *testing.T) {
	tests := []struct {
		name    string
		storage *MemStorage
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test checking maps for nil",
			storage: &MemStorage{
				mutex: &sync.RWMutex{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.storage.CheckStorage()
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() got error = %v, but wantErr %v", err, tt.wantErr)
			}

			require.NotNil(t, tt.storage.dataGauges)
			require.NotNil(t, tt.storage.dataCounters)
		})
	}
}

func TestMemStorage_String(t *testing.T) {
	store := NewMemStorage()
	ctx := context.Background()
	m1 := metrics.NewGaugeMetric("test1", 123.1)
	m2 := metrics.NewGaugeMetric("test2", 222.1)
	m3 := metrics.NewGaugeMetric("test3", 333.1)

	m4 := metrics.NewCounterMetric("test1", 111)
	m5 := metrics.NewCounterMetric("test2", 222)
	m6 := metrics.NewCounterMetric("test3", 333)

	err := store.AddMetric(ctx, m1)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m2)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m3)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m4)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m5)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m6)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	//fmt.Println(store.String(ctx))
	output1 := `[counters]
test1:111 | test2:222 | test3:333
[gauges]
test1:123.100000 | test2:222.100000 | test3:333.100000
`

	store2 := NewMemStorage()
	err = store2.AddMetric(ctx, m4)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store2.AddMetric(ctx, m5)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	fmt.Println(store2.String(ctx))
	output2 := `[counters]
test1:111 | test2:222`

	store3 := NewMemStorage()
	err = store3.AddMetric(ctx, m1)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store3.AddMetric(ctx, m2)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	fmt.Println(store3.String(ctx))
	output3 := `[gauges]
test1:123.100000 | test2:222.100000`

	emptyStore := NewMemStorage()
	tests := []struct {
		name    string
		storage *MemStorage
		want    string
	}{
		{
			name:    "Test full storage",
			storage: store,
			want:    output1,
		},
		{
			name:    "Test Only counters",
			storage: store2,
			want:    output2,
		},
		{
			name:    "Test Only gauges",
			storage: store3,
			want:    output3,
		},
		{
			name:    "Test empty storage",
			storage: emptyStore,
			want:    "empty storage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.storage.String(ctx); strings.TrimSpace(got) != strings.TrimSpace(tt.want) {
				t.Errorf("String() = --%v--, want --%v--", strings.TrimSpace(got), strings.TrimSpace(tt.want))
			}
		})
	}
}

func TestMemStorage_Close(t *testing.T) {
	ms := NewMemStorage()
	if err := ms.Close(context.Background()); err != nil {
		t.Errorf("Close() error = %v", err)
		return
	}
	require.Nil(t, ms.dataCounters)
	require.Nil(t, ms.dataGauges)

}

func TestMemStorage_UpdateMetric(t *testing.T) {
	store1 := NewMemStorage()
	ctx := context.Background()
	m := metrics.NewGaugeMetric("test1", 123.1)
	mNew := metrics.NewGaugeMetric("test1", 222.1)
	m1 := metrics.NewCounterMetric("test1", 111)
	m1New := metrics.NewCounterMetric("test1", 222)
	err := store1.AddMetric(ctx, m)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store1.AddMetric(ctx, m1)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}

	empty := NewMemStorage()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		store     *MemStorage
		oldMetric *metrics.Metrics
		newMetric *metrics.Metrics
	}{
		{
			name:      "Test update exist counter",
			store:     store1,
			newMetric: m1New,
		},
		{
			name:      "Test update exist gauge",
			store:     store1,
			newMetric: mNew,
		},
		{
			name:      "Test update not exist counter",
			store:     empty,
			newMetric: m1New,
		},
		{
			name:      "Test update not exist gauge",
			store:     empty,
			newMetric: mNew,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.newMetric.MType == metrics.Counter {
				oldDelta := tt.store.dataCounters[tt.newMetric.ID]
				if err := tt.store.UpdateMetric(ctx, tt.newMetric); err != nil {
					t.Errorf("UpdateMetric() got error = %v", err)
					return
				}
				updatedM := tt.store.dataCounters[tt.newMetric.ID]
				require.Equal(t, *tt.newMetric.Delta+oldDelta, updatedM)

			}

			if tt.newMetric.MType == metrics.Gauge {
				if err := tt.store.UpdateMetric(ctx, tt.newMetric); err != nil {
					t.Errorf("UpdateMetric() got error = %v", err)
					return
				}
				updatedM := tt.store.dataGauges[tt.newMetric.ID]
				require.Equal(t, *tt.newMetric.Value, updatedM)
			}

		})
	}
}

func TestMemStorage_AddMetric(t *testing.T) {
	store1 := NewMemStorage()
	ctx := context.Background()
	m := metrics.NewGaugeMetric("test1", 123.1)
	mNew := metrics.NewGaugeMetric("test1", 222.1)
	m1 := metrics.NewCounterMetric("test1", 111)
	m1New := metrics.NewCounterMetric("test1", 222)
	err := store1.AddMetric(ctx, m)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store1.AddMetric(ctx, m1)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}

	empty := NewMemStorage()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		store     *MemStorage
		oldMetric *metrics.Metrics
		newMetric *metrics.Metrics
	}{
		{
			name:      "Test update exist counter",
			store:     store1,
			newMetric: m1New,
		},
		{
			name:      "Test update exist gauge",
			store:     store1,
			newMetric: mNew,
		},
		{
			name:      "Test update not exist counter",
			store:     empty,
			newMetric: m1New,
		},
		{
			name:      "Test update not exist gauge",
			store:     empty,
			newMetric: mNew,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.newMetric.MType == metrics.Counter {
				oldDelta := tt.store.dataCounters[tt.newMetric.ID]
				if err := tt.store.AddMetric(ctx, tt.newMetric); err != nil {
					t.Errorf("UpdateMetric() got error = %v", err)
					return
				}
				updatedM := tt.store.dataCounters[tt.newMetric.ID]
				require.Equal(t, *tt.newMetric.Delta+oldDelta, updatedM)

			}

			if tt.newMetric.MType == metrics.Gauge {
				if err := tt.store.AddMetric(ctx, tt.newMetric); err != nil {
					t.Errorf("UpdateMetric() got error = %v", err)
					return
				}
				updatedM := tt.store.dataGauges[tt.newMetric.ID]
				require.Equal(t, *tt.newMetric.Value, updatedM)
			}

		})
	}
}

func conv(m *metrics.Metrics) string {
	var want string
	if m.MType == metrics.Counter {
		want = strconv.FormatInt(*m.Delta, 10)
	}
	if m.MType == metrics.Gauge {
		want = strconv.FormatFloat(*m.Value, 'f', -1, 64)
	}
	return want
}
func TestMemStorage_GetMetricString(t *testing.T) {
	store := NewMemStorage()
	ctx := context.Background()
	m1 := metrics.NewGaugeMetric("test1", 123.1)
	m2 := metrics.NewGaugeMetric("test2", 222.1)
	m3 := metrics.NewGaugeMetric("test3", 333.1)

	m4 := metrics.NewCounterMetric("test1", 111)
	m5 := metrics.NewCounterMetric("test2", 222)
	m6 := metrics.NewCounterMetric("test3", 333)

	err := store.AddMetric(ctx, m1)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m2)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m3)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m4)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m5)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m6)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}

	emptyStore := NewMemStorage()

	tests := []struct {
		name  string
		store *MemStorage
		m     *metrics.Metrics
		want  string
	}{
		{
			name:  "Test get counter exist",
			store: store,
			m:     m6,
			want:  conv(m6),
		},
		{
			name:  "Test get gauge exist",
			store: store,
			m:     m2,
			want:  conv(m2),
		},
		{
			name:  "Test get counter empty",
			store: emptyStore,
			m:     m6,
			want:  "",
		},
		{
			name:  "Test get gauge empty",
			store: emptyStore,
			m:     m2,
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := tt.store.GetMetricString(ctx, tt.m.MType, tt.m.ID)
			require.Equal(t, got, tt.want)
		})
	}
}

func TestMemStorage_GetMetric(t *testing.T) {
	store := NewMemStorage()
	ctx := context.Background()
	m1 := metrics.NewGaugeMetric("test1", 123.1)
	m2 := metrics.NewGaugeMetric("test2", 222.1)
	m3 := metrics.NewGaugeMetric("test3", 333.1)

	m4 := metrics.NewCounterMetric("test1", 111)
	m5 := metrics.NewCounterMetric("test2", 222)
	m6 := metrics.NewCounterMetric("test3", 333)

	err := store.AddMetric(ctx, m1)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m2)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m3)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m4)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m5)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m6)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}

	emptyStore := NewMemStorage()

	type args struct {
		store *MemStorage
		want  *metrics.Metrics
		key   string
		mType string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test get counter exist",
			args: args{
				store: store,
				want:  m6,
				key:   m6.ID,
				mType: m6.MType,
			},
			wantErr: false,
		},
		{
			name: "Test get gauge exist",
			args: args{
				store: store,
				want:  m2,
				key:   m2.ID,
				mType: m2.MType,
			},
			wantErr: false,
		},
		{
			name: "Test get counter empty",
			args: args{
				store: emptyStore,
				want:  nil,
				key:   m6.ID,
				mType: m6.MType,
			},
			wantErr: true,
		},
		{
			name: "Test get gauge empty",
			args: args{
				store: emptyStore,
				want:  nil,
				key:   m6.ID,
				mType: m6.MType,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.args.store.GetMetric(ctx, tt.args.key, tt.args.mType)
			if (err != nil) != tt.wantErr {
				t.Errorf("can't get metric got err: %v", err)
				return
			}
			if got != nil && tt.args.want != nil {
				require.Equal(t, got.ID, tt.args.want.ID)
				if got.MType == metrics.Counter {
					require.Equal(t, *tt.args.want.Delta, *got.Delta)
				}
				if got.MType == metrics.Gauge {
					require.Equal(t, *tt.args.want.Value, *got.Value)
				}
			}
		})
	}
}

func TestNewMemStorage(t *testing.T) {
	store := NewMemStorage()
	require.NotNil(t, store.dataGauges)
	require.NotNil(t, store.dataCounters)
	require.NotNil(t, store.mutex)
}

func TestMemStorage_GetAllMetrics(t *testing.T) {
	store := NewMemStorage()
	ctx := context.Background()
	m1 := metrics.NewGaugeMetric("test1", 123.1)
	m2 := metrics.NewGaugeMetric("test2", 222.1)
	m3 := metrics.NewGaugeMetric("test3", 333.1)

	m4 := metrics.NewCounterMetric("test4", 111)
	m5 := metrics.NewCounterMetric("test5", 222)
	m6 := metrics.NewCounterMetric("test6", 333)
	var allMetrics = make([]*metrics.Metrics, 6)
	allMetrics[0] = m1
	allMetrics[1] = m2
	allMetrics[2] = m3
	allMetrics[3] = m4
	allMetrics[4] = m5
	allMetrics[5] = m6

	err := store.AddMetric(ctx, m1)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m2)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m3)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m4)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m5)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}
	err = store.AddMetric(ctx, m6)
	if err != nil {
		fmt.Printf("can't add metric to storage: %v", err)
		return
	}

	emptyStore := NewMemStorage()
	allEmpty := make([]*metrics.Metrics, 0)

	type args struct {
		store *MemStorage
		slice []*metrics.Metrics
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test full storage",
			args: args{
				store: store,
				slice: allMetrics,
			},
		},
		{
			name: "Test empty storage",
			args: args{
				store: emptyStore,
				slice: allEmpty,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.store.GetAllMetrics(ctx)
			if err != nil {
				t.Errorf("GetAllMetrics() got error = %v", err)
				return
			}
			sort.Slice(got, func(i, j int) bool {
				return got[i].ID < got[i].ID
			})
			if !reflect.DeepEqual(got, tt.args.slice) {
				t.Errorf("GetAllMetrics() got:\n%v\nwant:\n%v", got, tt.args.slice)
			}
		})
	}
}

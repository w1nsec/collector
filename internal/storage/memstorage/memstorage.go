package memstorage

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/w1nsec/collector/internal/metrics"
)

var (
	errNotFound = fmt.Errorf("metric not found")
)

// for maps with counters and gauges
type MemStorage struct {
	dataCounters map[string]int64
	dataGauges   map[string]float64
	mutex        *sync.RWMutex
}

func (ms *MemStorage) Close(context.Context) error {
	ms.dataGauges = nil
	ms.dataCounters = nil
	return nil
}

func (ms *MemStorage) Init() error {
	ms.mutex.Lock()
	if ms.dataGauges == nil {
		ms.dataGauges = make(map[string]float64)
	}
	if ms.dataCounters == nil {
		ms.dataCounters = make(map[string]int64)
	}
	ms.mutex.Unlock()
	return nil
}

func (ms *MemStorage) CheckStorage() error {
	return ms.Init()
}

// Sting now changed to alphabetical order print
func (ms *MemStorage) String(ctx context.Context) string {
	ms.mutex.RLock()

	if len(ms.dataGauges) == 0 &&
		len(ms.dataCounters) == 0 {
		return "empty storage"
	}

	cSL := make([]string, len(ms.dataCounters))
	gSL := make([]string, len(ms.dataGauges))

	var s1 = "[counters]"
	var s2 = "[gauges]"
	i := 0
	for key, val := range ms.dataCounters {
		cSL[i] = fmt.Sprintf("%s:%d", key, val)
		i += 1
	}

	i = 0
	for key, val := range ms.dataGauges {
		gSL[i] = fmt.Sprintf("%s:%f", key, val)
		i += 1
	}
	sort.Strings(cSL)
	sort.Strings(gSL)

	s1 = fmt.Sprintf("%s\n%s\n",
		s1, strings.Join(cSL, " | "))
	s2 = fmt.Sprintf("%s\n%s\n",
		s2, strings.Join(gSL, " | "))

	var output string
	if len(ms.dataCounters) != 0 {
		output += s1
	}
	if len(ms.dataGauges) != 0 {
		output += s2
	}
	output += "\n"
	ms.mutex.RUnlock()

	return output

}

func (ms *MemStorage) UpdateMetric(ctx context.Context, newMetric *metrics.Metrics) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	switch newMetric.MType {
	case metrics.Gauge:
		ms.dataGauges[newMetric.ID] = *newMetric.Value
		return nil
	case metrics.Counter:
		ms.dataCounters[newMetric.ID] += *newMetric.Delta
		return nil
	}

	return nil
}

func (ms *MemStorage) AddMetric(ctx context.Context, newMetric *metrics.Metrics) error {
	return ms.UpdateMetric(ctx, newMetric)
}

func (ms *MemStorage) GetMetricString(ctx context.Context, mType, mName string) string {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()
	switch mType {
	case metrics.Gauge:
		val, ok := ms.dataGauges[mName]
		if !ok {
			return ""
		}
		return strconv.FormatFloat(val, 'f', -1, 64)
	case metrics.Counter:
		val, ok := ms.dataCounters[mName]
		if !ok {
			return ""
		}
		return strconv.FormatInt(val, 10)
	}
	return ""
}

func (ms *MemStorage) GetMetric(ctx context.Context, mName string, mType string) (*metrics.Metrics, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	switch mType {
	case metrics.Gauge:
		// for gauges
		for key, value := range ms.dataGauges {
			if key == mName {
				return &metrics.Metrics{
					ID:    key,
					MType: metrics.Gauge,
					Delta: nil,
					Value: &value,
				}, nil
			}
		}
	case metrics.Counter:
		// for counters
		for key, value := range ms.dataCounters {
			if key == mName {
				return &metrics.Metrics{
					ID:    key,
					MType: metrics.Counter,
					Delta: &value,
					Value: nil,
				}, nil
			}
		}
	}

	return nil, errNotFound
}

func NewMemStorage() *MemStorage {
	ms := &MemStorage{
		dataCounters: make(map[string]int64),
		dataGauges:   make(map[string]float64),
		mutex:        &sync.RWMutex{},
	}

	return ms
}

func (ms *MemStorage) GetAllMetrics(ctx context.Context) ([]*metrics.Metrics, error) {

	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	metricsSlice := make([]*metrics.Metrics, 0)

	// convert Gauges
	for key, val := range ms.dataGauges {
		// TODO what err, if I use pointer to val ???
		newVal := val
		/////

		metric := &metrics.Metrics{
			ID:    key,
			MType: metrics.Gauge,
			Value: &newVal,
		}

		metricsSlice = append(metricsSlice, metric)

	}

	// convert Counters
	for key, val := range ms.dataCounters {
		// TODO what err, if I use pointer to val ???
		newVal := val
		/////

		metric := &metrics.Metrics{
			ID:    key,
			MType: metrics.Counter,
			Delta: &newVal,
		}
		metricsSlice = append(metricsSlice, metric)
	}

	return metricsSlice, nil
}

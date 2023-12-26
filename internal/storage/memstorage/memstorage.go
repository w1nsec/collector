package memstorage

import (
	"context"
	"fmt"
	"strconv"
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

func (ms *MemStorage) String(ctx context.Context) string {
	ms.mutex.RLock()
	var s1 = "[counters]\n"
	var s2 = "[gauges]\n"
	count := 0
	numInLine := 5
	for key, val := range ms.dataCounters {
		s1 += fmt.Sprintf("%s:%d ", key, val)
		count++
		if count == numInLine {
			count = 0
			s1 += "\n"
		}
	}
	s1 += "\n"
	count = 0
	for key, val := range ms.dataGauges {
		s2 += fmt.Sprintf("%s:%f ", key, val)
		count++
		if count == numInLine {
			count = 0
			s2 += "\n"
		}
	}
	ms.mutex.RUnlock()
	s2 += "\n\n"
	return s1 + s2

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
	ms.UpdateMetric(ctx, newMetric)
	return nil
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

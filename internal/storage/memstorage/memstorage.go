package memstorage

import (
	"fmt"
	"github.com/w1nsec/collector/internal/metrics"
	"strconv"
)

type Storage interface {
	UpdateCounters(name string, value int64)
	UpdateGauges(name string, value float64)
	String() string
	GetMetricString(mType, mName string) string
	//GetOneMoreMetric(mType, mName string) string

	// valid
	GetMetric(mName string, mType string) *metrics.Metrics
	//UpdateMetrics(newMetrics []*metrics.Metrics) []error
	UpdateMetric(newMetric *metrics.Metrics)
	AddMetric(newMetric *metrics.Metrics)

	// add for increment9
	GetAllMetrics() []*metrics.Metrics
}

// for maps with counters and gauges
type MemStorage struct {
	dataCounters map[string]int64
	dataGauges   map[string]float64
	//metrics      map*metrics.Metrics
	//metrics []*metrics.Metrics
}

func (ms *MemStorage) String() string {
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
	s2 += "\n\n"

	return s1 + s2

}

func (ms *MemStorage) UpdateMetric(newMetric *metrics.Metrics) {
	switch newMetric.MType {
	case metrics.Gauge:
		ms.dataGauges[newMetric.ID] = *newMetric.Value
		return
	case metrics.Counter:
		ms.dataCounters[newMetric.ID] += *newMetric.Delta
		return
	}
}

func (ms *MemStorage) AddMetric(newMetric *metrics.Metrics) {
	ms.UpdateMetric(newMetric)
}

func (ms MemStorage) GetMetricString(mType, mName string) string {
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

func (ms MemStorage) GetMetric(mName string, mType string) *metrics.Metrics {
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
				}
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
				}
			}
		}
	}

	return nil
}

func (ms MemStorage) GetOneMetric(mName string) *metrics.Metrics {
	//for i, metric := range ms.metrics {
	//	if metric.MType == mName {
	//		return ms.metrics[i]
	//	}
	//}
	//return nil

	// for gauges
	for key, value := range ms.dataGauges {
		if key == mName {
			return &metrics.Metrics{
				ID:    key,
				MType: metrics.Gauge,
				Delta: nil,
				Value: &value,
			}
		}
	}

	// for counters
	for key, value := range ms.dataCounters {
		if key == mName {
			return &metrics.Metrics{
				ID:    key,
				MType: metrics.Counter,
				Delta: &value,
				Value: nil,
			}
		}
	}
	return nil
}

func NewMemStorage() *MemStorage {
	ms := new(MemStorage)
	ms.dataCounters = make(map[string]int64)
	ms.dataGauges = make(map[string]float64)
	return ms
}

func (ms *MemStorage) GetAllMetrics() []*metrics.Metrics {

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

	return metricsSlice
}

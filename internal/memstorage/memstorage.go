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
	GetMetric(mType, mName string) string
}

type MemStorage struct {
	dataCounters map[string]int64
	dataGauges   map[string]float64
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

func (ms *MemStorage) UpdateCounters(name string, value int64) {
	ms.dataCounters[name] += value
}

func (ms *MemStorage) UpdateGauges(name string, value float64) {
	ms.dataGauges[name] = value
}

func (ms MemStorage) GetMetric(mType, mName string) string {
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

func NewMemStorage() *MemStorage {
	ms := new(MemStorage)
	ms.dataCounters = make(map[string]int64)
	ms.dataGauges = make(map[string]float64)
	return ms
}

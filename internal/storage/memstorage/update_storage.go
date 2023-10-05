package memstorage

import (
	"fmt"
	"github.com/w1nsec/collector/internal/metrics"
)

func (ms *MemStorage) UpdateCounters(name string, value int64) {
	ms.dataCounters[name] += value
}

func (ms *MemStorage) UpdateGauges(name string, value float64) {
	ms.dataGauges[name] = value
}

func (ms *MemStorage) UpdateMetrics(newMetrics []*metrics.Metrics) []error {
	errors := make([]error, 0)
	for _, metric := range newMetrics {

		err := ms.updateMetric(metric)
		if err != nil {
			//log.Debug().Err(err)
			errors = append(errors, err)
		}
	}
	return errors
}

func (ms *MemStorage) updateMetric(metric *metrics.Metrics) error {
	val, err := reverseConvertOneMetric(metric)
	if err != nil {
		return err
	}

	if valInt, ok := val.(int64); ok {
		ms.UpdateCounters(metric.ID, valInt)
		return nil
	}

	if valFloat, ok := val.(float64); ok {
		ms.UpdateGauges(metric.ID, valFloat)
		return nil
	}

	return fmt.Errorf("can't convert metric to update")
}

func reverseConvertOneMetric(metric *metrics.Metrics) (interface{}, error) {
	switch metric.MType {
	case metrics.Gauge:
		return *metric.Value, nil
	case metrics.Counter:
		return *metric.Delta, nil
	}
	return nil, fmt.Errorf("wrong metric type ")
}

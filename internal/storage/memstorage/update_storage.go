package memstorage

import (
	"context"
	"fmt"
	"strings"

	"github.com/w1nsec/collector/internal/metrics"
)

func (ms *MemStorage) UpdateCounters(ctx context.Context, name string, value int64) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	ms.dataCounters[name] += value
	return nil
}

func (ms *MemStorage) UpdateGauges(ctx context.Context, name string, value float64) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	ms.dataGauges[name] = value
	return nil
}

func (ms *MemStorage) UpdateMetrics(ctx context.Context, newMetrics []*metrics.Metrics) error {

	errors := make([]string, 0)
	for _, metric := range newMetrics {
		err := ms.updateMetric(ctx, metric)
		if err != nil {
			//log.Debug().Err(err)
			errors = append(errors, err.Error())
		}
	}
	return fmt.Errorf(strings.Join(errors, " | "))
}

func (ms *MemStorage) updateMetric(ctx context.Context, metric *metrics.Metrics) error {

	val, err := reverseConvertOneMetric(metric)
	if err != nil {
		return err
	}

	if valInt, ok := val.(int64); ok {
		return ms.UpdateCounters(ctx, metric.ID, valInt)
	}

	if valFloat, ok := val.(float64); ok {
		return ms.UpdateGauges(ctx, metric.ID, valFloat)
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

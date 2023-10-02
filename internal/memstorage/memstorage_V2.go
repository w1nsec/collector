package memstorage

import (
	"encoding/json"
	"github.com/w1nsec/collector/internal/metrics"
)

func NewMetricStorage() *MetricStorage {
	ms := new(MetricStorage)
	ms.metrics = make([]*metrics.Metrics, 0)
	return ms
}

func (m *MetricStorage) UpdateCounters(name string, value int64) {
	//TODO implement me
	newMem := metrics.Metrics{
		ID:    name,
		MType: metrics.Counter,
		Delta: &value,
		Value: nil,
	}
	m.metrics = append(m.metrics, &newMem)
}

func (m *MetricStorage) UpdateGauges(name string, value float64) {
	//TODO implement me
	newMem := metrics.Metrics{
		ID:    name,
		MType: metrics.Counter,
		Delta: nil,
		Value: &value,
	}
	m.metrics = append(m.metrics, &newMem)
}

func (m MetricStorage) String() string {
	buf, err := json.Marshal(m.metrics)
	if err != nil {
		return ""
	}
	return string(buf)
}

func (m MetricStorage) GetMetric(mType, mName string) string {
	for id, metric := range m.metrics {
		if metric.ID == mName && metric.MType == mType {
			buf, err := json.Marshal(m.metrics[id])
			if err != nil {
				return err.Error()
			}
			return string(buf)
		}
	}
	return ""
}

func (m MetricStorage) GetOneMetric(mName string) *metrics.Metrics {
	for id, metric := range m.metrics {
		if metric.ID == mName {
			return m.metrics[id]
		}
	}
	return nil
}

func (m *MetricStorage) UpdateMetrics(newMetrics []*metrics.Metrics) []error {

	// TODO change this stupid and slow compare func
	for j, remoteMetric := range newMetrics {
		changed := false
		for i, localMetric := range m.metrics {
			if localMetric.ID == remoteMetric.ID {
				m.metrics[i] = newMetrics[j]
				changed = true
			}
		}
		if !changed {
			m.metrics = append(m.metrics, newMetrics[j])
		}

	}
	return nil
}

func (m *MetricStorage) UpdateMetric(metric *metrics.Metrics) {

	for i, locMetric := range m.metrics {
		if locMetric.ID == metric.ID {
			m.metrics[i] = metric
			return
		}
	}

	// metric not found yet
	m.metrics = append(m.metrics, metric)
}

func (m *MetricStorage) AddMetric(newMetric *metrics.Metrics) {
	m.metrics = append(m.metrics, newMetric)
}

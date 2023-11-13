package memstorage

import (
	"encoding/json"
	"github.com/w1nsec/collector/internal/metrics"
)

// Storage V2
type MetricStorage struct {
	metrics []*metrics.Metrics
}

func NewMetricStorage() *MetricStorage {
	ms := new(MetricStorage)
	ms.metrics = make([]*metrics.Metrics, 0)
	return ms
}

func (m *MetricStorage) UpdateCounters(name string, value int64) error {
	for ind, metric := range m.metrics {
		if metric.ID == name {
			*m.metrics[ind].Delta += value
		}
	}
	return nil
}

func (m *MetricStorage) UpdateGauges(name string, value float64) error {
	for ind, metric := range m.metrics {
		if metric.ID == name {
			*m.metrics[ind].Value = value
		}
	}
	return nil
}

func (m *MetricStorage) String() string {
	buf, err := json.Marshal(m.metrics)
	if err != nil {
		return ""
	}
	return string(buf)
}

func (m *MetricStorage) GetMetricString(mType, mName string) string {
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

func (m *MetricStorage) GetOneMetric(mName string) *metrics.Metrics {
	for id, metric := range m.metrics {
		if metric.ID == mName {
			return m.metrics[id]
		}
	}
	return nil
}

// No usage in interface
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

func (m *MetricStorage) UpdateMetric(newMetric *metrics.Metrics) error {

	for i, locMetric := range m.metrics {
		if locMetric.ID == newMetric.ID {
			m.metrics[i] = newMetric
			return nil
		}
	}

	// metric not found yet
	m.metrics = append(m.metrics, newMetric)
	return nil
}

func (m *MetricStorage) AddMetric(newMetric *metrics.Metrics) error {
	m.metrics = append(m.metrics, newMetric)
	return nil
}

func (m *MetricStorage) GetMetric(mName string, mType string) (*metrics.Metrics, error) {
	for _, metric := range m.metrics {
		if metric.ID == mName && metric.MType == mType {
			return metric, nil
		}
	}

	return nil, nil
}

func (m *MetricStorage) GetAllMetrics() ([]*metrics.Metrics, error) {
	return m.metrics, nil
}

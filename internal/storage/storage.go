package storage

import "github.com/w1nsec/collector/internal/metrics"

type Storage interface {
	UpdateCounters(name string, value int64) error
	UpdateGauges(name string, value float64) error
	String() string
	GetMetricString(mType, mName string) string

	// valid
	GetMetric(mName string, mType string) (*metrics.Metrics, error)

	UpdateMetric(newMetric *metrics.Metrics) error
	AddMetric(newMetric *metrics.Metrics) error

	// add for increment9 / increment3
	GetAllMetrics() ([]*metrics.Metrics, error)
}

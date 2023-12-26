package storage

import (
	"context"

	"github.com/w1nsec/collector/internal/metrics"
)

type Storage interface {
	UpdateCounters(ctx context.Context, name string, value int64) error
	UpdateGauges(ctx context.Context, name string, value float64) error
	String(ctx context.Context) string
	GetMetricString(ctx context.Context, mType, mName string) string

	// valid
	GetMetric(ctx context.Context, mName string, mType string) (*metrics.Metrics, error)

	UpdateMetric(ctx context.Context, newMetric *metrics.Metrics) error
	AddMetric(ctx context.Context, newMetric *metrics.Metrics) error

	// increment 12 many insert in DB
	UpdateMetrics(ctx context.Context, newMetrics []*metrics.Metrics) error

	// add for increment9 / increment3
	GetAllMetrics(ctx context.Context) ([]*metrics.Metrics, error)

	// for merging DBStorage and MemStorage
	Init() error
	CheckStorage() error

	// for shutdown
	Close(ctx context.Context) error
}

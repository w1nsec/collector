package metrics

import (
	"fmt"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

// Metrics is the main struct used for transport data
// between client and server
type Metrics struct {
	ID    string   `json:"id"`              // metric name
	MType string   `json:"type"`            // metric type, must be "gauge" or "counter" type
	Delta *int64   `json:"delta,omitempty"` // "counter" metric value (for "gauge" is nil)
	Value *float64 `json:"value,omitempty"` // "gauge" metric value (for "counter" is nil)
}

func (m Metrics) String() string {
	switch m.MType {
	case Gauge:
		return fmt.Sprintf("ID: %s | Type: %s | Value: %f", m.ID, m.MType, *m.Value)

	case Counter:
		return fmt.Sprintf("ID: %s | Type: %s | Value: %d", m.ID, m.MType, *m.Delta)

	}
	return fmt.Sprintf("ID: %s | Unsupported metric type", m.ID)
}

func NewCounterMetric(name string, value int64) *Metrics {
	metric := &Metrics{
		ID:    name,
		MType: Counter,
		Delta: &value,
	}
	return metric
}

func NewGaugeMetric(name string, value float64) *Metrics {

	metric := &Metrics{
		ID:    name,
		MType: Gauge,
		Value: &value,
	}
	//return metric, nil
	return metric
}

func Delete(metrics []*Metrics, ind int) []*Metrics {
	l := len(metrics)
	if ind >= l {
		return metrics
	}

	// delete last metric
	metrics[ind] = metrics[l-1]

	// TODO need this ?
	metrics[l-1] = nil

	return metrics[:l-1]
}

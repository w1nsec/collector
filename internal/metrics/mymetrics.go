package metrics

import (
	"fmt"
	"strconv"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

var (
	errConv = fmt.Errorf("error while converting \"mymetric\" to \"metric\"")
)

type MyMetrics struct {
	Value    string
	SendType string
}

func (m *MyMetrics) AddVal(n int) {
	// TODO parse error
	val, _ := strconv.ParseInt(m.Value, 10, 64)
	val += int64(n)
	m.Value = strconv.FormatInt(val, 10)
}

func ConvertMymetric2Metric(name string, mymetric MyMetrics) (*Metrics, error) {
	metric := &Metrics{}
	switch mymetric.SendType {
	case Gauge:
		val, err := strconv.ParseFloat(mymetric.Value, 64)
		if err != nil {
			return nil, err
		}
		metric = NewGaugeMetric(name, Gauge, val)
		if metric == nil {
			return nil, errConv
		}
	case Counter:
		val, err := strconv.Atoi(mymetric.Value)
		if err != nil {
			return nil, err
		}
		metric = NewCounterMetric(name, Counter, int64(val))
		if metric == nil {
			return nil, errConv
		}
	}
	return metric, nil
}

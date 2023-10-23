package metrics

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewCounterMetric(name, mType string, value int64) *Metrics {
	if mType == Gauge {
		return nil
	}
	metric := &Metrics{
		ID:    name,
		MType: mType,
		Delta: &value,
	}
	return metric
}

func NewGaugeMetric(name, mType string, value float64) *Metrics {
	if mType == Counter {
		//return nil, fmt.Errorf("wrong metric type, got: \"counter\", need: \"gauge\"")
		return nil
	}
	metric := &Metrics{
		ID:    name,
		MType: mType,
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

	// delete metric
	metrics[ind] = metrics[l-1]

	// TODO need this ?
	//newMetrics[len(newMetrics)-1] = nil

	return metrics[:l-1]
}

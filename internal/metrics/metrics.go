package metrics

import "strconv"

const (
	Gauge   = "gauge"
	Counter = "counter"
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

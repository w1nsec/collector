package handlers

import (
	"context"
	"fmt"
	"github.com/w1nsec/collector/internal/metrics"
	"net/http"
	"testing"
)

func (u *jsonGetUsecase) UpdateMetric(ctx context.Context, newMetric *metrics.Metrics) error {
	v, ok := u.m[newMetric.ID]
	// metric not exist, add it
	if !ok {
		u.m[newMetric.ID] = newMetric
		return nil
	}

	// metric already exist, update it
	switch v.MType {
	case metrics.Gauge:
		*u.m[newMetric.ID].Value = *newMetric.Value
		return nil
	case metrics.Counter:
		*u.m[newMetric.ID].Delta += *newMetric.Delta
		return nil
	}
	return fmt.Errorf("metric type not supported")
}

func TestJSONUpdateOneMetricHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		usecase updateOneMetricUsecase
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	// hint: starting github challenge 2024
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &JSONUpdateOneMetricHandler{
				usecase: tt.fields.usecase,
			}
			h.ServeHTTP(tt.args.w, tt.args.r)
		})
	}
}

package agent

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

// using storage in collection
func (agent Agent) CollectMetrics() {
	ctx := context.TODO()

	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	values := reflect.ValueOf(m)

	for _, name := range usedMemStats {
		structVal := values.FieldByName(name)
		var (
			val float64
		)

		if structVal.CanUint() {
			val = float64(structVal.Uint())
		}
		if structVal.CanInt() {
			val = float64(structVal.Int())
		}
		if structVal.CanFloat() {
			val = structVal.Float()
		}
		err := agent.store.UpdateMetric(ctx, &metrics.Metrics{
			ID:    name,
			MType: metrics.Gauge,
			Delta: nil,
			Value: &val,
		})
		if err != nil {
			log.Error().
				Err(err).Send()
			continue
		}
	}

	// Addition metrics
	// increase poll counter
	val := int64(1)
	err := agent.store.UpdateMetric(ctx, &metrics.Metrics{
		ID:    "PollCount",
		Delta: &val,
		MType: metrics.Counter,
	})
	if err != nil {
		log.Error().
			Err(err).Send()
	}
	// random value
	r := rand.New(rand.NewSource(time.Now().Unix()))
	randVal := float64(r.Intn(1000)) + r.Float64()
	err = agent.store.UpdateMetric(ctx, &metrics.Metrics{
		ID:    "RandomValue",
		Value: &randVal,
		MType: metrics.Gauge,
	})
	if err != nil {
		log.Error().
			Err(err).Send()
	}
}

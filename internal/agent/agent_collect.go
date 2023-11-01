package agent

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

// for storeV1 (maps with counter and gauge)
func (agent Agent) GetMetrics() {
	if agent.metrics == nil {
		agent.metrics = make(map[string]metrics.MyMetrics)
	}

	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	values := reflect.ValueOf(m)

	//metrics := make(map[string]interface{})

	// TODO Can you switch "reflect" to something other ???
	for _, name := range usedMemStats {
		if values.FieldByName(name).IsValid() {
			if values.FieldByName(name).CanInt() {
				agent.metrics[name] = metrics.MyMetrics{
					Value:    strconv.FormatInt(values.FieldByName(name).Int(), 10),
					SendType: metrics.Gauge,
				}
				continue
			}
			if values.FieldByName(name).CanUint() {
				agent.metrics[name] = metrics.MyMetrics{
					Value:    strconv.FormatUint(values.FieldByName(name).Uint(), 10),
					SendType: metrics.Gauge,
				}
				continue
			}
			if values.FieldByName(name).CanFloat() {
				agent.metrics[name] = metrics.MyMetrics{
					Value:    strconv.FormatFloat(values.FieldByName(name).Float(), 'f', -1, 64),
					SendType: metrics.Gauge,
				}
				continue
			}
		}

	}

	// Addition metrics
	metric, ok := agent.metrics["PollCount"]
	if !ok {
		agent.metrics["PollCount"] = metrics.MyMetrics{
			Value:    "1",
			SendType: metrics.Counter,
		}
	} else {
		metric.AddVal(1)
		agent.metrics["PollCount"] = metric
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	randVal := float64(r.Intn(1000)) + r.Float64()
	agent.metrics["RandomValue"] = metrics.MyMetrics{
		Value:    strconv.FormatFloat(randVal, 'f', -1, 64),
		SendType: metrics.Gauge,
	}
}

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

package agent

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/w1nsec/collector/internal/metrics"
)

// UsedMemStats - params that should collect agent by task
var UsedMemStats = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

// CollectMainMetrics collect UsedMemStats metrics only
func (agent Agent) CollectMainMetrics(ctx context.Context) {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	values := reflect.ValueOf(m)

	for _, name := range UsedMemStats {
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

// CollectGopsutilMetrics collect addition metrics^
// - CPUutilization
// - FreeMemory
// - TotalMemory
// added by task in increment15 / iter15 gopsutil
func (agent Agent) CollectGopsutilMetrics(ctx context.Context) {

	// gather memory metrics
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Error().
			Err(err).Send()
		return
	}

	// TotalMemory
	totalMem := float64(v.Total)
	err = agent.store.UpdateMetric(ctx, &metrics.Metrics{
		ID:    "TotalMemory",
		Value: &totalMem,
		MType: metrics.Gauge,
	})
	if err != nil {
		log.Error().
			Err(err).Send()
	}

	// FreeMemory
	freeMem := float64(v.Free)
	err = agent.store.UpdateMetric(ctx, &metrics.Metrics{
		ID:    "FreeMemory",
		Value: &freeMem,
		MType: metrics.Gauge,
	})
	if err != nil {
		log.Error().
			Err(err).Send()
	}

	// CPU utilization
	utilizations, err := cpu.Percent(time.Millisecond, true)
	if err != nil {
		fmt.Println(err)
		log.Error().
			Err(err).Send()
	}
	metricName := "CPUutilization"
	for ind, utiliz := range utilizations {
		err = agent.store.UpdateMetric(ctx, &metrics.Metrics{
			ID:    fmt.Sprintf("%s%d", metricName, ind),
			Value: &utiliz,
			MType: metrics.Gauge,
		})
		if err != nil {
			log.Error().
				Err(err).Send()
		}
	}
}

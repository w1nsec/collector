package main

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"
)

var usedMemStats = []string{
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

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	runtime.GOMAXPROCS(3)
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)
	metrics := make(map[string]interface{})
	//working := make(chan bool)

	mut := sync.Mutex{}
	for {
		select {
		case t1 := <-pollTicker.C:
			// because we can appeal to "metrics" at the same time !!!
			fmt.Println("Receiving:", t1.Format(time.TimeOnly))
			go func() {
				mut.Lock()
				GetMetrics(metrics)
				mut.Unlock()
			}()

			// Variant #1
			//go func(c chan bool) {
			//	c <- true
			//	GetMetrics(metrics)
			//	<-c
			//}(working)
		case t2 := <-reportTicker.C:
			fmt.Println("- Sending:", t2.Format(time.TimeOnly))
			go func() {
				mut.Lock()
				SendMetrics(metrics)
				mut.Unlock()
			}()

			// Variant #1
			//go func() {
			//	working <- true
			//	SendMetrics(metrics)
			//	<-working
			//}()

		}
	}

}

func SendMetrics(metrics map[string]interface{}) {
	//fmt.Println("Sending:", time.Now().Format(time.TimeOnly))
}

func GetMetrics(metrics map[string]interface{}) {
	//fmt.Println("Receiving:", time.Now().Format(time.TimeOnly))
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	values := reflect.ValueOf(m)

	//metrics := make(map[string]interface{})
	for _, name := range usedMemStats {
		if values.FieldByName(name).IsValid() {
			if values.FieldByName(name).CanInt() {
				metrics[name] = values.FieldByName(name).Int()
			}
			if values.FieldByName(name).CanUint() {
				metrics[name] = values.FieldByName(name).Uint()
			}
			if values.FieldByName(name).CanFloat() {
				metrics[name] = values.FieldByName(name).Float()
			}
		}
	}
}

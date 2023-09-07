package main

import (
	"fmt"
	"runtime"
)

// Alloc
// BuckHashSys
// Frees
// GCCPUFraction
// GCSys
// HeapAlloc
// HeapIdle
// HeapInuse
// HeapObjects
// HeapReleased
// HeapSys
// LastGC
// Lookups
// MCacheInuse
// MCacheSys
// MSpanInuse
// MSpanSys
// Mallocs
// NextGC
// NumForcedGC
// NumGC
// OtherSys
// PauseTotalNs
// StackInuse
// StackSys
// Sys
// TotalAlloc
func main() {

	// get descriptions for all supported metrics
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)

	fmt.Println(m)
	m.BySize()
}

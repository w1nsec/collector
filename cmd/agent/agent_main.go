package main

import (
	"fmt"
	"github.com/w1nsec/collector/internal/agent"
	"runtime"
	"sync"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	runtime.GOMAXPROCS(3)
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)
	//metrics := make(map[string]interface{})
	addr := "localhost:8080"

	mAgent, err := agent.NewAgent(addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	mAgent.GetMetrics()
	mAgent.SendMetrics()
	return

	mAgent.GetMetrics()
	// Variant #2
	mut := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for t1 := range pollTicker.C {
			mut.Lock()
			fmt.Println("Receiving:", t1.Format(time.TimeOnly))
			mAgent.GetMetrics()
			mut.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		for t2 := range reportTicker.C {
			mut.Lock()
			fmt.Println("- Sending:", t2.Format(time.TimeOnly))
			mAgent.SendMetrics()
			mut.Unlock()
		}
	}()
	wg.Wait()
	// ======================================== //

	// Variant #1
	// because we can appeal to "metrics" at the same time !!!
	// ======================================== //
	//
	//  working := make(chan bool)
	//	for {
	//		select {
	//		case t1 := <-pollTicker.C:
	//			fmt.Println("Receiving:", t1.Format(time.TimeOnly))
	//			go func(c chan bool) {
	//				c <- true
	//				GetMetrics(metrics)
	//				<-c
	//			}(working)
	//		case t2 := <-reportTicker.C:
	//			fmt.Println("- Sending:", t2.Format(time.TimeOnly))
	//			go func() {
	//				working <- true
	//				SendMetrics(metrics)
	//				<-working
	//			}()
	//		}
	//	}
	// ======================================== //
}

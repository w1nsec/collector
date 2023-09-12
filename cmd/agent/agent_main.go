package main

import (
	"flag"
	"fmt"
	"github.com/w1nsec/collector/internal/agent"
	"time"
)

const (
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
)

func main() {
	//runtime.GOMAXPROCS(3)

	//metrics := make(map[string]interface{})
	//addr := "localhost:8080"
	var (
		addr                         string
		pollInterval, reportInterval int
	)
	flag.StringVar(&addr, "a", "localhost:8080",
		"address for metric server")

	flag.IntVar(&pollInterval, "r", int(defaultPollInterval.Seconds()),
		"frequency of gathering metrics")
	flag.IntVar(&reportInterval, "p", int(defaultReportInterval.Seconds()),
		"frequency of sending metrics")
	flag.Parse()

	mAgent, err := agent.NewAgent(addr, pollInterval, reportInterval)
	if err != nil {
		fmt.Println(err)
		return
	}
	mAgent.Start()

	// Variant #2
	// ======================================== //
	//mut := sync.Mutex{}
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//go func() {
	//	for t1 := range pollTicker.C {
	//		mut.Lock()
	//		fmt.Println("Receiving:", t1.Format(time.TimeOnly))
	//		mAgent.GetMetrics()
	//		mut.Unlock()
	//	}
	//}()
	//
	//wg.Add(1)
	//go func() {
	//	for t2 := range reportTicker.C {
	//		mut.Lock()
	//		fmt.Println("- Sending:", t2.Format(time.TimeOnly))
	//		mAgent.SendMetrics()
	//		mut.Unlock()
	//	}
	//}()
	//wg.Wait()
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

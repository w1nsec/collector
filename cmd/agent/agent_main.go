package main

import (
	"fmt"
	"github.com/w1nsec/collector/internal/agent"
	"github.com/w1nsec/collector/internal/config"
	"log"
)

func main() {

	//metrics := make(map[string]interface{})
	//addr := "localhost:8080"
	var (
		addr                         string
		pollInterval, reportInterval int
	)

	config.AgentSelectArgs(&addr, &pollInterval, &reportInterval)

	mAgent, err := agent.NewAgent(addr, pollInterval, reportInterval)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Fatal(mAgent.Start())

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

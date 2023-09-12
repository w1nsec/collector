package main

import (
	"flag"
	"fmt"
	"github.com/w1nsec/collector/internal/agent"
	"os"
	"strconv"
	"time"
)

const (
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
)

func selectArgs(addr *string, pollInterval, reportInterval *int) {
	var (
		flagAddr          string
		flagPoll, flagRep int
	)
	flag.StringVar(&flagAddr, "a", "localhost:8080",
		"address for metric server")
	flag.IntVar(&flagPoll, "p", int(defaultPollInterval.Seconds()),
		"frequency of gathering metrics")
	flag.IntVar(&flagRep, "r", int(defaultReportInterval.Seconds()),
		"frequency of sending metrics")
	flag.Parse()

	if *addr = os.Getenv("ADDRESS"); *addr == "" {
		*addr = flagAddr
	}

	envPoll, err := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	if err == nil {
		*pollInterval = envPoll
	} else {
		*pollInterval = flagPoll
	}

	envRep, err := strconv.Atoi(os.Getenv("REPORT_INTERVAL"))
	if err == nil {
		*reportInterval = envRep
	} else {
		*reportInterval = flagRep
	}

}

func main() {
	//runtime.GOMAXPROCS(3)

	//metrics := make(map[string]interface{})
	//addr := "localhost:8080"
	var (
		addr                         string
		pollInterval, reportInterval int
	)

	selectArgs(&addr, &pollInterval, &reportInterval)

	fmt.Println(addr, pollInterval, reportInterval)
	return
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

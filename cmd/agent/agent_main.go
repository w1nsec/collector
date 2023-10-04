package main

import (
	"fmt"
	"github.com/w1nsec/collector/internal/agent"
	"github.com/w1nsec/collector/internal/config"
	"log"
)

func main() {
	
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

}

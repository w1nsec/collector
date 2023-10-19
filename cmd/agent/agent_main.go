package main

import (
	"fmt"
	"github.com/w1nsec/collector/internal/agent"
	agent2 "github.com/w1nsec/collector/internal/config/agent"
	"log"
)

func main() {

	var (
		addr                         string
		pollInterval, reportInterval int
	)

	agent2.AgentSelectArgs(&addr, &pollInterval, &reportInterval)

	mAgent, err := agent.NewAgent(addr, pollInterval, reportInterval)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Fatal(mAgent.Start())

}

package main

import (
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/agent"
	config "github.com/w1nsec/collector/internal/config/agent"
	"github.com/w1nsec/collector/internal/logger"
)

func main() {

	var (
		addr                         string
		pollInterval, reportInterval int
		loggerLevel                  = "info"
	)

	config.AgentSelectArgs(&addr, &pollInterval, &reportInterval)
	err := logger.Initialize(loggerLevel)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	mAgent, err := agent.NewAgent(addr, pollInterval, reportInterval)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	err = mAgent.Start()
	log.Fatal().Err(err).Send()

}

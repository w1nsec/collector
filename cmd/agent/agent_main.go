package main

import (
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/agent"
	config "github.com/w1nsec/collector/internal/config/agent"
	"github.com/w1nsec/collector/internal/logger"
)

func main() {

	var (
		args config.Args
	)

	config.AgentSelectArgs(&args)
	err := logger.Initialize(args.LogLevel)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	mAgent, err := agent.NewAgent(args)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	err = mAgent.Start()
	log.Fatal().Err(err).Send()

}

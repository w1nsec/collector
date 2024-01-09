package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/app/agent"
	config "github.com/w1nsec/collector/internal/config/agent"
	"github.com/w1nsec/collector/internal/logger"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	args := config.AgentSelectArgs()
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
	err = mAgent.Start(ctx)
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}
	defer mAgent.Close()
	log.Info().Msg("Closing agent app: successful")
}

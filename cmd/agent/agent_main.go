package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/agent"
	config "github.com/w1nsec/collector/internal/config/agent"
	"github.com/w1nsec/collector/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var (
		args config.Args
	)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
	defer cancel()

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
	err = mAgent.Start(ctx)
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}
	log.Info().Msg("Closing agent app: successful")
}

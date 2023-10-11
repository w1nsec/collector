package main

import (
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config"
	"github.com/w1nsec/collector/internal/service"
)

func main() {

	var args config.Args

	config.ServerArgsParse(&args)
	log.Info().
		Str("addr", args.Addr).
		Str("log", args.LogLevel).Send()

	Service, err := service.NewService(args)
	if err != nil {
		log.Fatal().Err(err).Send()

	}

	log.Fatal().Err(Service.Start()).Send()
}

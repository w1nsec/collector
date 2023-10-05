package main

import (
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config"
	"github.com/w1nsec/collector/internal/server"
)

func main() {

	var args config.Args

	config.ServerArgsParse(&args)

	log.Info().
		Str("addr", args.Addr).
		Str("log", args.LogLevel).Send()

	srv, err := server.NewServer(args)
	if err != nil {
		log.Fatal().Err(err).Send()

	}
	//srv.AddMux(mux)

	log.Fatal().Err(srv.Start()).Send()
}

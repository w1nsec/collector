package main

import (
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config"
	"github.com/w1nsec/collector/internal/server"
)

func main() {
	var (
		addr     string
		logLevel string
	)
	config.ServerArgsParse(&addr, &logLevel)

	log.Info().
		Str("addr", addr).
		Str("log", logLevel).Send()

	srv, err := server.NewServer(addr, logLevel)
	if err != nil {
		log.Fatal().Err(err).Send()

	}
	//srv.AddMux(mux)

	log.Fatal().Err(srv.Start()).Send()
}

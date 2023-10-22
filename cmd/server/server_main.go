package main

import (
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/app"
)

func main() {
	serverApp, err := app.NewAppServer()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Fatal().Err(serverApp.Run()).Send()
}

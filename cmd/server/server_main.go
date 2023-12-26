package main

import (
	"context"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/app/server"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
	defer stop()

	serverApp, err := server.NewAppServer()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Fatal().Err(serverApp.Run(ctx)).Send()
}

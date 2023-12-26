package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/app/server"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
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

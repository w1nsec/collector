package main

import (
	"fmt"
	"github.com/w1nsec/collector/internal/config"
	"github.com/w1nsec/collector/internal/server"
	"log"
)

func main() {
	var (
		addr     string
		logLevel string
	)
	config.ServerArgsParse(&addr, &logLevel)

	fmt.Println(addr, logLevel)
	srv, err := server.NewServer(addr, logLevel)
	if err != nil {
		log.Fatal(err)
	}
	//srv.AddMux(mux)

	log.Fatal(srv.Start())
}

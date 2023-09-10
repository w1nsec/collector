package main

import (
	"github.com/w1nsec/collector/internal/server"
	"log"
)

func main() {

	addr := "localhost:8080"

	srv, err := server.NewMetricServer(addr)
	if err != nil {
		log.Fatal(err)
	}
	//srv.AddMux(mux)

	log.Fatal(srv.Start())
}

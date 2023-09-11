package main

import (
	"flag"
	"github.com/w1nsec/collector/internal/server"
	"log"
)

func main() {

	addr := flag.String("a", "localhost:8080", "address for server")
	flag.Parse()

	srv, err := server.NewMetricServer(*addr)
	if err != nil {
		log.Fatal(err)
	}
	//srv.AddMux(mux)

	log.Fatal(srv.Start())
}

package main

import (
	"flag"
	"fmt"
	"github.com/w1nsec/collector/internal/server"
	"log"
)

var (
	addr string
)

func initFlags() {
	flag.StringVar(&addr, "a", "localhost:8080", "address for server")
}

func main() {

	initFlags()
	flag.Parse()
	fmt.Println(addr)
	srv, err := server.NewMetricServer(addr)
	if err != nil {
		log.Fatal(err)
	}
	//srv.AddMux(mux)

	log.Fatal(srv.Start())
}

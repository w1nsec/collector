package main

import (
	"flag"
	"fmt"
	"github.com/w1nsec/collector/internal/server"
	"log"
	"os"
)

func selectArgs(addr *string) {
	*addr = os.Getenv("ADDRESS")

	var flagAddr string
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address for server")
	flag.Parse()

	if *addr == "" {
		*addr = flagAddr
	}

}

func main() {
	var (
		addr string
	)
	selectArgs(&addr)

	fmt.Println(addr)
	srv, err := server.NewMetricServer(addr)
	if err != nil {
		log.Fatal(err)
	}
	//srv.AddMux(mux)

	log.Fatal(srv.Start())
}

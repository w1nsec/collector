package main

import (
	"github.com/w1nsec/collector/internal/handlers"
	"github.com/w1nsec/collector/internal/memstorage"
	"github.com/w1nsec/collector/internal/server"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	addr := "localhost:8080"

	store := memstorage.NewMemStorage()
	srv, err := server.NewMetricServer(addr, store, mux)
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/update/", handlers.UpdateMetricsHandle(srv))

	log.Fatal(http.ListenAndServe(addr, mux))
}

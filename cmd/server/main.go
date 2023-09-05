package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	addr := "localhost:8080"

	log.Fatal(http.ListenAndServe(addr, mux))
}

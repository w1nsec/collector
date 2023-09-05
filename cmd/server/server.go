package main

import (
	"github.com/w1nsec/collector/internal/memstorage"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var store *memstorage.MemStorage

func main() {
	mux := http.NewServeMux()
	store = memstorage.NewMemStorage()

	addr := "localhost:8080"
	mux.HandleFunc("/update/", updateMetricsHandle)

	log.Fatal(http.ListenAndServe(addr, mux))
}

func updateMetricsHandle(rw http.ResponseWriter, r *http.Request) {
	metricsInfo, found := strings.CutPrefix(r.URL.Path, "/update/")
	if !found {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	pieces := strings.Split(metricsInfo, "/")
	log.Println(pieces)
	if len(pieces) != 3 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	mType := pieces[0]
	mName := pieces[1]
	mValue := pieces[2]

	switch mType {
	case "gauge":
		val, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		store.UpdateGauges(mName, val)
		log.Println(store)
	case "counter":
		val, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		store.UpdateCounters(mName, val)
		log.Println(store)
	default:
		rw.WriteHeader(http.StatusBadRequest)
	}

	rw.WriteHeader(http.StatusOK)
}

package handlers

import (
	"github.com/w1nsec/collector/internal/server"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func UpdateMetricsHandle(srv *server.MetricServer) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
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
			srv.Store.UpdateGauges(mName, val)
			log.Println(srv.Store)
		case "counter":
			val, err := strconv.ParseInt(mValue, 10, 64)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			srv.Store.UpdateCounters(mName, val)
			log.Println(srv.Store)
		default:
			rw.WriteHeader(http.StatusBadRequest)
		}

		rw.WriteHeader(http.StatusOK)
		return
	}
}

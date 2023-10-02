package handlers

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/memstorage"
	"github.com/w1nsec/collector/internal/metrics"
	"io"
	"net/http"
)

//type JSONUpdateHandler struct {
//	store memstorage.Storage
//}

// func (h JSONUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

// JSON format update
func JSONUpdateHandler(store memstorage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Debug().Err(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		recMetrics := make([]*metrics.Metrics, 0)
		err = json.Unmarshal(body, &recMetrics)
		if err != nil {
			log.Debug().Err(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		errs := store.UpdateMetrics(recMetrics)
		if errs != nil {
			for _, err = range errs {
				log.Debug().Err(err)
			}
			w.WriteHeader(http.StatusInternalServerError)
		}

		// All is good, return updated metrics
		// Debug version
		// TODO change request metrics by name to:  body-request -> body-response resent
		ids := make([]string, 0)
		for _, metric := range recMetrics {
			if metric != nil {
				ids = append(ids, metric.ID)
			}
		}

		respMetrics := make([]*metrics.Metrics, 0)
		for i := 0; i < len(ids); i++ {
			metric := store.GetOneMetric(ids[i])
			if metric != nil {
				respMetrics = append(respMetrics, metric)
			}
		}

		bodyResponse, err := json.Marshal(respMetrics)
		if err != nil {
			log.Debug().Err(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// write answer
		_, err = w.Write(bodyResponse)
		if err != nil {
			log.Debug().Err(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

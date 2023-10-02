package handlers

import (
	"bytes"
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
		// TODO change request metrics by name to body-request -> body-response
		ids := make([]string, 0)
		types := make([]string, 0)
		for _, metric := range recMetrics {
			if metric != nil {
				ids = append(ids, metric.ID)
				types = append(types, metric.MType)
			}
		}

		// How to do it???
		if len(types) != len(ids) {
			log.Fatal().Msgf("I can't do it")
		}
		bodyResponse := make([]byte, 0)
		buffer := bytes.NewBuffer(bodyResponse)
		for i := 0; i < len(types); i++ {
			mStr := store.GetMetric(types[i], ids[i])
			buffer.WriteString(mStr)
		}

		// write answer
		w.Write(buffer.Bytes())
	}
}

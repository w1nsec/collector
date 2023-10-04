package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/memstorage"
	"github.com/w1nsec/collector/internal/metrics"
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
			log.Err(err).Send()
			//log.Info().Msgf("err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		reqMetrica := metrics.Metrics{}
		err = json.Unmarshal(body, &reqMetrica)
		if err != nil {
			log.Err(err).Send()
			//log.Info().Msgf("err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		errs := store.UpdateMetrics([]*metrics.Metrics{&reqMetrica})
		if len(errs) != 0 {
			for _, err = range errs {
				log.Err(err)
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// All is good, return updated metrics
		// Debug version
		// TODO change request metrics by name to:  body-request -> body-response resent
		// ids := make([]string, 0)
		// for _, metric := range recMetrics {
		// 	if metric != nil {
		// 		ids = append(ids, metric.ID)
		// 	}
		// }

		// respMetrics := make([]*metrics.Metrics, 0)
		// for i := 0; i < len(ids); i++ {
		// 	metric := store.GetOneMetric(ids[i])
		// 	if metric != nil {
		// 		respMetrics = append(respMetrics, metric)
		// 	}
		// }

		val := store.GetMetric(reqMetrica.MType, reqMetrica.ID)
		if val == "" {
			// метрика не обнаружена
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}
		respMetrica := store.GetOneMetric(reqMetrica.ID)

		bodyResponse, err := json.Marshal(&respMetrica) // todo все так-ти брать из хранилища
		if err != nil {
			log.Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// write answer
		_, err = w.Write(bodyResponse)
		if err != nil {
			log.Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func JSONValueHandler(store memstorage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Err(err).Send()
			//log.Info().Msgf("err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		reqMetrica := metrics.Metrics{}
		err = json.Unmarshal(body, &reqMetrica)
		if err != nil {
			log.Err(err).Send()
			//log.Info().Msgf("err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		val := store.GetMetric(reqMetrica.MType, reqMetrica.ID)
		if val == "" {
			// метрика не обнаружена
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}
		respMetrica := store.GetOneMetric(reqMetrica.ID)

		bodyResponse, err := json.Marshal(&respMetrica)
		if err != nil {
			log.Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// write answer
		_, err = w.Write(bodyResponse)
		if err != nil {
			log.Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// это можно не писать, по дефолту статусОК
		//w.WriteHeader(http.StatusOK)
	}
}

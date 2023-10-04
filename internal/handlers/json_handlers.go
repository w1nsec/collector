package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/memstorage"
	"github.com/w1nsec/collector/internal/metrics"
	"io"
	"net/http"
	"strings"
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

		recMetrics := make([]*metrics.Metrics, 0)
		err = json.Unmarshal(body, &recMetrics)
		if err != nil {
			log.Err(err).Send()
			//log.Info().Msgf("err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		errs := store.UpdateMetrics(recMetrics)
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

func JSONUpdateOneMetricHandler(store memstorage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			log.Error().
				Err(fmt.Errorf("wrong method for %s", r.URL.RawPath)).
				Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		valid := false
		vals := r.Header.Values("content-type")
		for _, val := range vals {
			if val == "application/json" {
				valid = true
				break
			}
		}
		if !valid {
			log.Error().
				Err(fmt.Errorf("invalid \"content-type\": %s", strings.Join(vals, ";"))).
				Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var metric *metrics.Metrics
		metric = new(metrics.Metrics)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error().
				Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &metric)
		if err != nil {
			log.Error().
				Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Info().
			RawJSON("metric", body).
			Msg("Request")

		// Check, that metric contains values
		if (metric.Delta == nil && metric.Value == nil) ||
			metric.ID == "" {
			log.Error().
				Err(fmt.Errorf("metric doesn't contain any value")).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		store.UpdateMetric(metric)

		// Debug version
		// TODO change request metrics by name to:  body-request -> body-response resent (when debug done)

		metric = store.GetMetric(metric.ID, metric.MType)
		if metric == nil {
			log.Error().
				Err(fmt.Errorf("metric \"%s\" not found in store", metric.ID)).Send()
			w.WriteHeader(http.StatusNotFound)
			return
		}

		body, err = json.Marshal(metric)
		if err != nil {
			log.Error().
				Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info().
			RawJSON("metric", body).
			Msg("Response")

		w.Header().Set("content-type", "application/json")
		_, err = w.Write(body)
		if err != nil {
			log.Error().
				Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func JSONGetMetricHandler(store memstorage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			log.Error().
				Err(fmt.Errorf("wrong method for %s", r.URL.RawPath)).
				Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		valid := false
		vals := r.Header.Values("content-type")
		for _, val := range vals {
			if val == "application/json" {
				valid = true
				break
			}
		}
		if !valid {
			log.Error().
				Err(fmt.Errorf("invalid \"content-type\": %s", strings.Join(vals, ";"))).
				Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var metric *metrics.Metrics
		metric = new(metrics.Metrics)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error().
				Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &metric)
		if err != nil {
			log.Error().
				Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Info().
			RawJSON("metric", body).
			Msg("Request")

		// Debug version
		// TODO change request metrics by name to:  body-request -> body-response resent (when debug done)

		// Check, that metric contains values
		if metric.ID == "" {
			log.Error().
				Err(fmt.Errorf("metric doesn't contain ID")).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		metric = store.GetMetric(metric.ID, metric.MType)
		if metric == nil {
			log.Error().
				Err(fmt.Errorf("metric \"%s\" not found in store", metric.ID)).Send()
			w.WriteHeader(http.StatusNotFound)
			return
		}

		body, err = json.Marshal(metric)
		if err != nil {
			log.Error().
				Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info().
			RawJSON("metric", body).
			Msg("Response")

		w.Header().Set("content-type", "application/json")
		_, err = w.Write(body)
		if err != nil {
			log.Error().
				Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

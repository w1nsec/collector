package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage"
	"io"
	"net/http"
	"strings"
)

func JSONUpdateOneMetricHandler(store storage.Storage) func(w http.ResponseWriter, r *http.Request) {
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
			if strings.Contains(val, "application/json") {
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

		var metric = new(metrics.Metrics)

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

		if metric == nil {
			log.Error().
				Err(fmt.Errorf("metric is nil")).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Check, that metric contains values
		if (metric.Delta == nil && metric.Value == nil) ||
			metric.ID == "" {
			log.Error().
				Err(fmt.Errorf("metric doesn't contain any value")).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = store.UpdateMetric(metric)
		if err != nil {
			log.Error().
				Err(err).
				Msg("can't update metric")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Debug version
		// TODO change request metrics by name to:  body-request -> body-response resent (when debug done)

		retMetric, _ := store.GetMetric(metric.ID, metric.MType)
		if retMetric == nil {
			log.Error().
				Err(fmt.Errorf("metric \"%s\" not found in store",
					metric.ID)).Send()
			w.WriteHeader(http.StatusNotFound)
			return
		}

		body, err = json.Marshal(retMetric)
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

func JSONGetMetricHandler(store storage.Storage) func(w http.ResponseWriter, r *http.Request) {
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
			if strings.Contains(val, "application/json") {
				valid = true
				break
			}
		}
		if !valid {
			log.Error().
				Err(fmt.Errorf("invalid \"content-type\": %s",
					strings.Join(vals, ";"))).
				Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var metric = new(metrics.Metrics)

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
			Str("method", r.Method).
			Str("url", r.URL.RawPath).
			Msg("Request")

		// Debug version
		// TODO change request metrics by name to:  body-request -> body-response resent (when debug done)

		if metric == nil {
			log.Error().
				Err(fmt.Errorf("metric is nil")).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if metric.ID == "" {
			log.Error().
				Err(fmt.Errorf("metric doesn't contain ID")).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Check, that metric contains values
		if metric.ID == "" {
			log.Error().
				Err(fmt.Errorf("metric doesn't contain ID")).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		retMetric, _ := store.GetMetric(metric.ID, metric.MType)
		if retMetric == nil {
			log.Error().
				Err(fmt.Errorf("metric \"%s\" not found in store",
					metric.ID)).Send()
			w.WriteHeader(http.StatusNotFound)
			return
		}

		body, err = json.Marshal(retMetric)
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

// increment 12
func JSONUpdateMetricsHandler(store storage.Storage) func(w http.ResponseWriter, r *http.Request) {
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
			if strings.Contains(val, "application/json") {
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

		var newMetrics = make([]*metrics.Metrics, 0)

		err := json.NewDecoder(r.Body).Decode(&newMetrics)
		if err != nil {
			log.Error().
				Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		//log.Info().
		//	RawJSON("metric", body).
		//	Msg("Request")

		// Check, that metric contains values
		errors := make([]string, 0)
		for _, m := range newMetrics {
			if (m.Delta == nil && m.Value == nil) ||
				m.ID == "" {
				err := fmt.Errorf("metric \"%s\"doesn't contain any value", m.ID)
				log.Error().
					Err(err).Send()
				//w.WriteHeader(http.StatusInternalServerError)
				//return
				errors = append(errors, err.Error())
				continue
			}

			err = store.UpdateMetric(m)
			if err != nil {
				log.Error().
					Err(err).
					Msgf("can't update metric \"%s\"", m.ID)
				err = fmt.Errorf("can't update metric \"%s\"", m.ID)
				errors = append(errors, err.Error())
				continue
			}

		}

		// If no one metric is saved
		if len(errors) == len(newMetrics) {
			io.WriteString(w, "Don't save any metric")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		_, err = io.WriteString(w, strings.Join(errors, " | "))
		if err != nil {
			log.Error().
				Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
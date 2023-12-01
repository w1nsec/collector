package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/service"
	"io"
	"net/http"
	"strings"
)

func JSONUpdateOneMetricHandler(store *service.MetricService) func(w http.ResponseWriter, r *http.Request) {
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

		err = store.UpdateMetric(r.Context(), metric)
		if err != nil {
			log.Error().
				Err(err).
				Msg("can't update metric")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Debug version
		// TODO change request metrics by name to:  body-request -> body-response resent (when debug done)

		retMetric, _ := store.GetMetric(r.Context(), metric.ID, metric.MType)
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

func JSONGetMetricHandler(service *service.MetricService) func(w http.ResponseWriter, r *http.Request) {
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
			Str("url", r.URL.RequestURI()).
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

		retMetric, _ := service.GetMetric(r.Context(), metric.ID, metric.MType)
		if retMetric == nil {
			log.Error().
				Err(fmt.Errorf("metric \"%s\" not found in service",
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
func JSONUpdateMetricsHandler(service *service.MetricService) func(w http.ResponseWriter, r *http.Request) {
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

		// check repeats
		mNames := make(map[string]string, 0)

		// Check, that metric contains values
		errors := make([]string, 0)
		for ind, m := range newMetrics {
			if (m.Delta == nil && m.Value == nil) ||
				m.ID == "" {
				err := fmt.Errorf("metric \"%s\"doesn't contain any value", m.ID)
				log.Error().
					Err(err).Send()

				// delete wrong metric
				newMetrics = metrics.Delete(newMetrics, ind)

				errors = append(errors, err.Error())
				continue
			}
			mNames[m.ID] = m.MType

		}

		// log localerrors
		if len(errors) != 0 {
			log.Error().
				Err(fmt.Errorf(strings.Join(errors, " | "))).
				Send()
		}

		err = service.UpdateMetrics(r.Context(), newMetrics)
		if err != nil {
			log.Error().
				Err(err).
				Send()
			io.WriteString(w, "Don't save any metric")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get updated metrics
		updatedMetrics := make([]*metrics.Metrics, 0)
		for mName, mType := range mNames {
			metric, err := service.GetMetric(r.Context(), mName, mType)
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}
			updatedMetrics = append(updatedMetrics, metric)
		}

		w.Header().Set("content-type", "application/json")
		err = json.NewEncoder(w).Encode(updatedMetrics)
		if err != nil {
			log.Error().
				Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

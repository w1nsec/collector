package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
)

type jsonUpdateUsecase interface {
	jsonGetMetricUsecase
	UpdateMetrics(ctx context.Context, newMetrics []*metrics.Metrics) error
}

type JSONUpdateMetricsHandler struct {
	usecase jsonUpdateUsecase
}

func NewJSONUpdateMetricsHandler(usecase jsonUpdateUsecase) *JSONUpdateMetricsHandler {
	return &JSONUpdateMetricsHandler{usecase: usecase}
}

// increment 12
func (h *JSONUpdateMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

	err = h.usecase.UpdateMetrics(r.Context(), newMetrics)
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
		metric, err := h.usecase.GetMetric(r.Context(), mName, mType)
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

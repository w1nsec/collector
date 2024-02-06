package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
)

type updateOneMetricUsecase interface {
	jsonGetMetricUsecase
	UpdateMetric(ctx context.Context, newMetric *metrics.Metrics) error
}

type JSONUpdateOneMetricHandler struct {
	usecase updateOneMetricUsecase
}

func NewJSONUpdateOneMetricHandler(usecase updateOneMetricUsecase) *JSONUpdateOneMetricHandler {
	return &JSONUpdateOneMetricHandler{usecase: usecase}
}

// JSONUpdateOneMetricHandler godoc
// @Tags JSON Metrics
// @Summary Update one metric in json
// @Description Update one metric in json by it name and value (in json)
// @Accept application/json
// @Produce application/json
// @Param bodyReq body metrics.Metrics true "JSON, contains Metric ID,type,value"
// @Success 200 {object} metrics.Metrics "OK, update/add metric value"
// @Success 404 {string} string "Not found, can't found metric by ID"
// @Failure 500 {string} string "Internal error"
// @Router /update/ [post]
func (h *JSONUpdateOneMetricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		log.Error().
			Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info().
		Str("metric", metric.ID).
		Msg("Request")

	// Check, that metric contains values
	if (metric.Delta == nil && metric.Value == nil) ||
		metric.ID == "" {
		log.Error().
			Err(fmt.Errorf("metric doesn't contain any value")).Send()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.usecase.UpdateMetric(r.Context(), metric)
	if err != nil {
		log.Error().
			Err(err).
			Msg("can't update metric")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	retMetric, err := h.usecase.GetMetric(r.Context(), metric.ID, metric.MType)
	if retMetric == nil || err != nil {
		log.Error().
			Err(fmt.Errorf("metric \"%s\" not found in store",
				metric.ID)).Send()
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(retMetric)
	if err != nil {
		log.Error().
			Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info().
		Str("metric", metric.ID).
		Msg("Response")

	w.WriteHeader(http.StatusOK)
}

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

type jsonGetMetricUsecase interface {
	GetMetric(ctx context.Context, mName string, mType string) (*metrics.Metrics, error)
}

type JSONGetMetricHandler struct {
	usecase jsonGetMetricUsecase
}

func NewJSONGetMetricHandler(usecase jsonGetMetricUsecase) *JSONGetMetricHandler {
	return &JSONGetMetricHandler{usecase: usecase}
}

// JSONGetMetricHandler godoc
// @Tags JSON Metrics
// @Summary Get one metric in json
// @Description Get one metric in json by it name (in json)
// @Accept application/json
// @Produce application/json
// @Param bodyReq body metrics.Metrics true "MetricID"
// @Success 200 {object} metrics.Metrics "OK, return metric value"
// @Success 404 {string} string "Not found, can't found metric by ID"
// @Failure 500 {string} string "Internal error"
// @Router /value/ [post]
func (h *JSONGetMetricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

	// Check, that metric contains ID
	if metric.ID == "" {
		log.Error().
			Err(fmt.Errorf("metric doesn't contain ID")).Send()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check, that metric contains values
	//if metric.Value == nil && metric.Delta == nil {
	//	log.Error().
	//		Err(fmt.Errorf("metric doesn't contain any data")).Send()
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	retMetric, err := h.usecase.GetMetric(r.Context(), metric.ID, metric.MType)
	if retMetric == nil || err != nil {
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

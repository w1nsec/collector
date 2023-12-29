package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type updateGaugeUsecase interface {
	UpdateGauges(ctx context.Context, name string, value float64) error
}

type UpdateGaugeHandler struct {
	//counterUsecase updateCountersUsecase
	gaugeUsecase updateGaugeUsecase
}

func NewUpdateGaugeHandler(gaugeUsecase updateGaugeUsecase) *UpdateGaugeHandler {
	return &UpdateGaugeHandler{gaugeUsecase: gaugeUsecase}
}

// UpdateGaugeHandler godoc
// @Tags Update Metrics
// @Summary Update Gauges
// @Description  Update gauge metric
// @ID updateGauge
// @Produce text/plain
// @Param name path string true "Metric name"
// @Param value path string true "Metric value"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Wrong request, can't parse number"
// @Failure 500 {string} string "Internal error"
// @Router /update/gauge/{name}/{value} [post]
func (h *UpdateGaugeHandler) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.gaugeUsecase.UpdateGauges(r.Context(), name, val)
	if err != nil {
		wr.WriteHeader(http.StatusInternalServerError)
		log.Error().
			Err(err).Msg("can't update gauges")
		return
	}
	wr.WriteHeader(http.StatusOK)
}

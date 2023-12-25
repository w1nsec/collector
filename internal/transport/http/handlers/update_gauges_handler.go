package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
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

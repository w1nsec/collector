package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type updateCountersUsecase interface {
	UpdateCounters(ctx context.Context, name string, value int64) error
}

type UpdateCountersHandler struct {
	counterUsecase updateCountersUsecase
}

func NewUpdateCountersHandler(counterUsecase updateCountersUsecase) *UpdateCountersHandler {
	return &UpdateCountersHandler{counterUsecase: counterUsecase}
}

// UpdateCountersHandler godoc
// @Tags Update Metrics
// @Summary Update Counters
// @Description  Update Counter metric
// @ID updateCounter
// @Produce text/plain
// @Param name path string true "Metric name"
// @Param value path string true "Metric value"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Wrong request, can't parse number"
// @Failure 500 {string} string "Internal error"
// @Router /update/counter/{name}/{value} [post]
func (h *UpdateCountersHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.counterUsecase.UpdateCounters(r.Context(), name, val)
	if err != nil {
		log.Error().Err(err).Msg("can't update counters")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

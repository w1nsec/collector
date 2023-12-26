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

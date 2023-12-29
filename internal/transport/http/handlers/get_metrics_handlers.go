package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type getMetricUsecase interface {
	GetMetricString(ctx context.Context, mType string, mName string) string
}

type GetMetricHandler struct {
	usecase getMetricUsecase
}

func NewGetMetricHandler(usecase getMetricUsecase) *GetMetricHandler {
	return &GetMetricHandler{usecase: usecase}
}

func (h *GetMetricHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")

	value := h.usecase.GetMetricString(r.Context(), mType, mName)

	// if metric not found
	if value == "" {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(NotFound))
		return
	}

	rw.Header().Add("Content-type", "text/plain")
	_, err := rw.Write([]byte(value))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)

}

type getAllMetricsUsecase interface {
	String(ctx context.Context) string
}

type GetMetricsHandler struct {
	usecase getAllMetricsUsecase
}

func NewGetMetricsHandler(usecase getAllMetricsUsecase) *GetMetricsHandler {
	return &GetMetricsHandler{usecase: usecase}
}

func (h *GetMetricsHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	page := fmt.Sprintf("<html><body>%s</body></html>",
		h.usecase.String(r.Context()))

	rw.Header().Set("content-type", "text/html")

	_, err := fmt.Fprint(rw, page)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Error().
			Err(err).Msg("can't output storage in html")
		return
	}

}

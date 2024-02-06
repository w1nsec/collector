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

// GetMetricHandler godoc
// @Tags Get Metrics
// @Summary Get one metric
// @Description Get one metric by name
// @Produce text/plain
// @Param mType path string true "Metric type"
// @Param mName path string true "Metric name"
// @Success 200 {string} string "OK, return metric value"
// @Failure 400 {string} string "Wrong request, can't parse number"
// @Failure 404 {string} string "Metric not found"
// @Failure 500 {string} string "Internal error, can't write response body"
// @Router /{mType}/{mName} [get]
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
		_, err := rw.Write([]byte(NotFound))
		if err != nil {
			log.Error().Err(err).
				Msg("can't write body to ResponseWriter")
		}
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

// GetMetricsHandler godoc
// @Tags Get Metrics
// @Summary Get All Metric
// @Description Get All Metric in html format
// @Produce text/html
// @Success 200 {string} string "OK, return metric value"
// @Failure 500 {string} string "Internal error, can't write response body"
// @Router / [get]
func (h *GetMetricsHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	page := fmt.Sprintf("<html><body>%s</body></html>",
		h.usecase.String(r.Context()))

	rw.Header().Set("content-type", "text/html")

	defer fmt.Fprint(rw, page)
}

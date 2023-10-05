package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/memstorage"
	"net/http"
)

func GetMetric(store memstorage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "mType")
		mName := chi.URLParam(r, "mName")

		value := store.GetMetricString(mType, mName)
		// metric not found
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
}

func GetAllMetrics(store memstorage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		page := fmt.Sprintf("<html><body>%s</body></html>", store.String())
		rw.Header().Set("content-type", "text/html")
		_, err := fmt.Fprint(rw, page)
		if err != nil {
			log.Error().Err(err).Send()
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

package handlers

import (
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/service"
	"net/http"
)

func CheckDBConnectionHandler(service *service.MetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := service.CheckStorage(r.Context())
		if err != nil {
			http.Error(w, "Can't connect to DB", http.StatusInternalServerError)
			log.Error().Err(err).Send()
			return
		}
		_, err = w.Write([]byte("DB available"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Can't write slice to body", http.StatusInternalServerError)
			log.Error().Err(err).Send()
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}

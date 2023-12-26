package handlers

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

type checkStorageUsecase interface {
	CheckStorage() error
}

type CheckDBConnectionHandler struct {
	checkUsecase checkStorageUsecase
}

func NewCheckDBConnectionHandler(checkUsecase checkStorageUsecase) *CheckDBConnectionHandler {
	return &CheckDBConnectionHandler{checkUsecase: checkUsecase}
}

func (h *CheckDBConnectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := h.checkUsecase.CheckStorage()
	if err != nil {
		http.Error(w, "Can't connect to DB", http.StatusInternalServerError)
		log.Error().
			Err(err).
			Msgf("can't connect to DB")
		return
	}
	_, err = w.Write([]byte("DB available"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Error().
			Err(err).
			Msg("can't connect to DB")
		return
	}

	w.WriteHeader(http.StatusOK)

}

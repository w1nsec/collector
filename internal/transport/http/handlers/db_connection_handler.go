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

// CheckDBConnectionHandler godoc
// @Tags DB Check
// @Tags.description "Check connection to DB"
// @Summary Check DB
// @Description Check connection to DB
// @ID checkDB
// @Produce text/plain
// @Success 200 {string} string "DB available"
// @Failure 500 {string} string "Internal error"
// @Router /ping/ [get]
func (h *CheckDBConnectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := h.checkUsecase.CheckStorage()
	if err != nil {
		http.Error(w, "Can't connect to DB", http.StatusInternalServerError)
		log.Error().
			Err(err).
			Msgf("can't connect to DB")
		return
	}
	defer w.Write([]byte("DB available"))
}

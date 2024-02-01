package handlers

import (
	"io"
	"net/http"
)

// Pong godoc
// @Tags Ping debug
// @Summary Test handler for service
// @Description Request echo reply from server
// @ID pongHandler
// @Produce text/plain
// @Success 200 {object} string "pong"
// @Failure 500 {string} string "Internal error"
// @Router /echoping [get]
func Pong(w http.ResponseWriter, r *http.Request) {
	defer io.WriteString(w, "pong\n")
}

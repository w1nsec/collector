package handlers

import (
	"io"
	"net/http"
)

func Pong(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "pong\n")
	if err != nil {
		http.Redirect(w, r, "error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

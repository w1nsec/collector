package handlers

import (
	"io"
	"net/http"
)

func pong(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "pong\n")
	if err != nil {
		http.Redirect(w, r, "error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

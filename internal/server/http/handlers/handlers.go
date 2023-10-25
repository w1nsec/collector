package handlers

import (
	"net/http"
)

var (
	NotFound = "404 page not found\n"
)

func NotFoundHandle(rw http.ResponseWriter, r *http.Request) {
	http.NotFound(rw, r)
}

func BadRequest(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}

package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/w1nsec/collector/internal/memstorage"
	"net/http"
)

var (
	NotFound = "404 page not found\n"
)

func GetMetric(store memstorage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "mType")
		mName := chi.URLParam(r, "mName")

		value := store.GetMetric(mType, mName)
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

func NotFoundHandle(rw http.ResponseWriter, r *http.Request) {
	http.NotFound(rw, r)
}

func BadRequest(rw http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	rw.WriteHeader(http.StatusBadRequest)
}

package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func GzipDecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		encoding := r.Header.Get("content-encoding")
		if encoding == "" {
			next.ServeHTTP(w, r)
			return
		}
		if encoding != "gzip" {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().
				Err(fmt.Errorf("usupported content-encoding: %s", encoding)).Send()
			return
		}

		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().
				Err(fmt.Errorf("usupported content-encoding: %s", encoding)).Send()
			return
		}

		r.Body = reader
		next.ServeHTTP(w, r)
	})
}

package middlewares

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type responseData struct {
	respSize   int
	statusCode int
}
type loggerRW struct {
	http.ResponseWriter
	*responseData
}

func (lrw *loggerRW) Write(buf []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(buf)
	lrw.responseData.respSize = size
	return size, err
}

func (lrw *loggerRW) WriteHeader(statusCode int) {
	lrw.ResponseWriter.WriteHeader(statusCode)
	lrw.responseData.statusCode = statusCode
}

func LoggingMiddleware(h http.Handler) http.Handler {
	logFunc := func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		uri := r.RequestURI
		start := time.Now()

		var (
			data = responseData{}
			lrw  = loggerRW{ResponseWriter: w, responseData: &data}
		)
		h.ServeHTTP(&lrw, r)
		duration := time.Since(start)

		log.Info().
			Msgf("Request: %s %s time: %fs",
				method, uri, duration.Seconds())

		log.Info().
			Int("status", lrw.statusCode).
			Int("resp size", lrw.responseData.respSize).
			Msg("Response:")
	}
	return http.HandlerFunc(logFunc)
}

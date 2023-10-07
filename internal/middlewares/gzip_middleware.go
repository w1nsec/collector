package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"compress/gzip"
	"github.com/rs/zerolog/log"
)

type gzipResponseWriter struct {
	http.ResponseWriter
}

// TODO review this func
func (gRW gzipResponseWriter) Write(data []byte) (int, error) {
	// should response be compressed ??
	// checking content-type
	shouldCompress := false
	vals := gRW.Header().Values("content-type")
	for _, val := range vals {
		if strings.Contains(val, "application/json") ||
			strings.Contains(val, "text/html") {
			shouldCompress = true
			break
		}
	}
	// say WRONG to compress
	if !shouldCompress {
		log.Debug().
			Err(fmt.Errorf("content-type is wrong")).
			Msg("response should NOT be compressed")

		return gRW.ResponseWriter.Write(data)
	}

	gzWriter, err := gzip.NewWriterLevel(gRW.ResponseWriter, gzip.BestSpeed)
	if err != nil {
		log.Error().
			Err(err).Send()
		http.Error(gRW.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return 0, err
	}
	defer gzWriter.Close()
	gRW.Header().Set("Content-Encoding", "gzip")
	log.Debug().
		Msg("response SHOULD be compressed")

	return gzWriter.Write(data)
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Reading request

		// checking current encoding
		encoded := false
		vals := r.Header.Values("Content-encoding")
		for _, val := range vals {
			if strings.Contains(val, "gzip") {
				encoded = true
				break
			}
		}

		// text encoded, so, read now from gzip Reader
		if encoded {
			// decompress
			// create gzip reader
			gzReader, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Error().
					Err(err).Send()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//defer r.Body.Close()
			r.Body = gzReader
			defer gzReader.Close()
			log.Debug().
				Msg("Client have compressed BODY")
		}

		// --------------------------------
		// checking accept-encoding
		acceptEnc := false
		vals = r.Header.Values("accept-encoding")
		for _, val := range vals {
			if strings.Contains(val, "gzip") {
				acceptEnc = true
				break
			}
		}

		// if client do not accept compression ...
		if !acceptEnc {
			log.Debug().
				Msg("Client don't accept compression")
			next.ServeHTTP(w, r)
			return
		}

		// prepare for compression
		next.ServeHTTP(gzipResponseWriter{ResponseWriter: w}, r)
	})
}

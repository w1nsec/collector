package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/utils/signing"
)

var (
	errWrongSigning = fmt.Errorf("wrong signing")
)

type signingMidl struct {
	secret     string
	hmacHeader string
}

func NewSigningMidl(secret string) *signingMidl {
	return &signingMidl{secret: secret, hmacHeader: "HashSHA256"}
}

func (s signingMidl) Signing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if secret not setup, then just pass through
		if s.secret == "" {
			next.ServeHTTP(w, r)
			return
		}

		// if request doesn't contain need header
		headerVal := r.Header.Get(s.hmacHeader)
		if headerVal == "" {
			next.ServeHTTP(w, r)
			return
		}

		// check signing
		buf := bytes.NewBuffer(nil)
		_, err := io.Copy(buf, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Send()
			return
		}

		if !signing.CheckSigning(buf.Bytes(), []byte(headerVal), []byte(s.secret)) {
			w.WriteHeader(http.StatusBadRequest)
			log.Error().Err(errWrongSigning).Send()
			return
		}

		newBody := io.NopCloser(buf)
		r.Body = newBody
		next.ServeHTTP(w, r)

	})
}

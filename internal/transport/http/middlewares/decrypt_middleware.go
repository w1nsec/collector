package middlewares

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config"
	"github.com/w1nsec/go-examples/crypto"
)

type DecryptMiddleware struct {
	privKey *rsa.PrivateKey
}

func NewDecryptMiddleware(key *rsa.PrivateKey) *DecryptMiddleware {
	return &DecryptMiddleware{
		privKey: key,
	}
}

// Handle - provide decrypt function for message body
// 1. Get encrypted AES key from header
// 2. Decrypt it with RSA PrivateKey
// 3. Decrypt body with AES key
func (h *DecryptMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// no need for decrypting
		if h.privKey == nil {
			next.ServeHTTP(rw, r)
			return
		}

		// only for requests with body
		if r.Body != nil {

			// get encrypted aes from header
			header := r.Header.Get(config.CryptoHeader)
			// header not set, no encryption, next
			if header == "" {
				next.ServeHTTP(rw, r)
				return
			}

			encAES, err := base64.StdEncoding.DecodeString(header)
			if err != nil {
				log.Error().Err(err).Send()
				http.Error(rw, "can't decrypt request", http.StatusInternalServerError)
				return
			}

			// decrypt aes key
			key, err := rsa.DecryptPKCS1v15(rand.Reader, h.privKey, encAES)
			if err != nil {
				log.Error().Err(err).Send()
				http.Error(rw, "can't decrypt request", http.StatusInternalServerError)
				return
			}

			// get encrypted body
			cData, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error().Err(err).Send()
				http.Error(rw, "can't decrypt request", http.StatusInternalServerError)
				return
			}

			data, err := crypto.DecryptAES(cData, key)
			if err != nil {
				log.Error().Err(err).Send()
				http.Error(rw, "can't decrypt request", http.StatusInternalServerError)
				return
			}

			//log.Info().Msg(string(data))

			// set decrypted body back to request and other handlers
			buf := bytes.NewBuffer(data)
			newBody := io.NopCloser(buf)
			r.Body = newBody
			next.ServeHTTP(rw, r)
		}
	})
}

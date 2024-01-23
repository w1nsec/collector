package signing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/rs/zerolog/log"
)

func CreateSigning(data, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	_, err := h.Write(data)
	if err != nil {
		log.Error().Err(err).Send()
		return nil
	}
	sign := h.Sum(nil)

	return []byte(base64.StdEncoding.EncodeToString(sign))
}

func CheckSigning(data, sign, key []byte) bool {
	genSign := CreateSigning(data, key)
	return string(genSign) == string(sign)
}

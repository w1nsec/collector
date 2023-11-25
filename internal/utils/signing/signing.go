package signing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func CreateSigning(data, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	sign := h.Sum(nil)

	return []byte(base64.StdEncoding.EncodeToString(sign))
}

func CheckSigning(data, sign, key []byte) bool {
	genSign := CreateSigning(data, key)
	return string(genSign) == string(sign)
}

package middlewares

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w1nsec/collector/internal/config"
	"github.com/w1nsec/go-examples/crypto"
)

func TestNewDecryptHandler(t *testing.T) {
	type args struct {
		key *rsa.PrivateKey
	}
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		return
	}

	handler := NewDecryptMiddleware(key)
	require.NotNil(t, handler, "return nil from constructor")
	require.Equal(t, *key, *handler.privKey)
}

func TestDecryptMiddleware_Handle(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		return
	}

	aesKey := crypto.FillAESKey([]byte("supersecretkey"))
	pub := key.Public().(*rsa.PublicKey)
	aesEnc1, err := crypto.EncryptRSA(aesKey, pub)
	if err != nil {
		fmt.Println(err)
		return
	}
	aesEnc1Str := base64.StdEncoding.EncodeToString(aesEnc1)

	body := []byte(strings.Repeat("keep calm and write code", 10))
	bodyEnc, err := crypto.EncryptAES(body, aesKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	key2, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		return
	}

	type args struct {
		privKey    *rsa.PrivateKey
		method     string
		body       []byte
		headers    map[string]string
		wantStatus int
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
		{
			name: "Test no priv key",
			args: args{
				privKey:    nil,
				method:     http.MethodGet,
				body:       nil,
				wantStatus: http.StatusOK,
			},
		},
		{
			name: "Test no body",
			args: args{
				privKey:    key,
				method:     http.MethodGet,
				body:       nil,
				wantStatus: http.StatusOK,
			},
		},
		{
			name: "Test NO base64 encoded key",
			args: args{
				privKey: key,
				method:  http.MethodPost,
				body:    []byte("standart body"),
				headers: map[string]string{
					config.CryptoHeader: "*sd*",
				},
				wantStatus: http.StatusInternalServerError,
			},
		},
		{
			name: "Test can't decrypt AES",
			args: args{
				privKey: key2,
				method:  http.MethodPost,
				body:    []byte("standart body"),
				headers: map[string]string{
					config.CryptoHeader: aesEnc1Str,
				},
				wantStatus: http.StatusInternalServerError,
			},
		},
		{
			name: "Test can't decrypt body",
			args: args{
				privKey: key,
				method:  http.MethodPost,
				body:    []byte("standart body"),
				headers: map[string]string{
					config.CryptoHeader: aesEnc1Str,
				},
				wantStatus: http.StatusInternalServerError,
			},
		},
		{
			name: "Test all is good",
			args: args{
				privKey: key,
				method:  http.MethodPost,
				body:    bodyEnc,
				headers: map[string]string{
					config.CryptoHeader: aesEnc1Str,
				},
				wantStatus: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DecryptMiddleware{
				privKey: tt.args.privKey,
			}

			body := bytes.NewBuffer(tt.args.body)
			req := httptest.NewRequest(tt.args.method, "http://localhost:8000", body)
			for k, v := range tt.args.headers {
				req.Header.Add(k, v)
			}
			res := httptest.NewRecorder()
			handler(res, req)

			got := h.Handle(handler)
			got.ServeHTTP(res, req)
			result := res.Result()
			defer result.Body.Close()

			require.Equal(t, tt.args.wantStatus, result.StatusCode,
				fmt.Sprintf("wrong status code = %v, want %v",
					result.StatusCode, tt.args.wantStatus))

		})
	}
}

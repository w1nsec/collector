package middlewares

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
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

	handler := NewDecryptHandler(key)
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

	type args struct {
		privKey    *rsa.PrivateKey
		method     string
		body       []byte
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
			name: "Test no body",
			args: args{
				privKey:    key,
				method:     http.MethodPost,
				body:       []byte("standart body"),
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

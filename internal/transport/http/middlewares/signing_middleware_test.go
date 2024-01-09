package middlewares

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"github.com/w1nsec/collector/internal/utils/signing"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_signingMidl_Signing(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	body := strings.Repeat("new body", 20)
	secret := "supersecret"
	hmacHeader := "HashSHA256"
	signing := signing.CreateSigning([]byte(body), []byte(secret))

	type args struct {
		method  string
		secret  string
		hmacVal string
		body    []byte
	}
	tests := []struct {
		name  string
		args  args
		valid bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Valid Signing",
			args: args{
				method:  http.MethodPost,
				secret:  secret,
				hmacVal: string(signing),
				body:    []byte(body),
			},
			valid: true,
		},
		{
			name: "Test Invalid Signing",
			args: args{
				method:  http.MethodPost,
				secret:  secret,
				hmacVal: "invalid sign",
				body:    []byte(body),
			},
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &signingMidl{
				secret:     tt.args.secret,
				hmacHeader: hmacHeader,
			}

			body := bytes.NewBuffer(tt.args.body)
			req := httptest.NewRequest(tt.args.method, "http://localhost:8000/", body)
			req.Header.Add(hmacHeader, tt.args.hmacVal)
			recoder := httptest.NewRecorder()

			handler(recoder, req)

			got := s.Signing(handler)
			got.ServeHTTP(recoder, req)

			res := recoder.Result()
			defer res.Body.Close()

			if !tt.valid {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				return
			}

			require.Equal(t, http.StatusOK, res.StatusCode)
		})
	}
}

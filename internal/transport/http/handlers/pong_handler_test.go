package handlers

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPong(t *testing.T) {
	//handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	var pingResp = "pong\n"

	type args struct {
		method string
		body   []byte
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "Test GET ping",
			args: args{
				method: http.MethodGet,
				body:   nil,
			},
		},
		{
			name: "Test POST ping",
			args: args{
				method: http.MethodPost,
				body:   []byte("ping"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rw := httptest.NewRecorder()
			body := bytes.NewBuffer(tt.args.body)
			r := httptest.NewRequest(tt.args.method, "http://localhost:8000/ping", body)
			Pong(rw, r)

			result := rw.Result()
			defer result.Body.Close()

			resBody, err := io.ReadAll(result.Body)
			if err != nil {
				t.Errorf("can't read body from response: %v", err)
				return
			}

			require.Equal(t, pingResp, string(resBody))

		})
	}
}

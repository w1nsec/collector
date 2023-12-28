package middlewares

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/w1nsec/collector/internal/utils/compression/gzip"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGzipDecompressMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	text := strings.Repeat("simpletext", 10)
	validBody, err := gzip.Compress([]byte(text))
	if err != nil {
		fmt.Println("can't generate gzip body")
		return
	}

	type args struct {
		method         string
		compressHeader string
		body           []byte
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		// TODO: Add test cases.
		{
			name: "Test Valid GET request",
			args: args{
				method:         http.MethodGet,
				compressHeader: "gzip",
				body:           nil,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Test Valid gzip compressed message",
			args: args{
				method:         http.MethodPost,
				compressHeader: "gzip",
				body:           validBody,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Test Invalid gzip compressed message",
			args: args{
				method:         http.MethodPost,
				compressHeader: "gzip",
				body:           []byte{},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "Test Invalid encoding",
			args: args{
				method:         http.MethodPost,
				compressHeader: "unzip",
				body:           []byte{},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBuffer(tt.args.body)
			req := httptest.NewRequest(tt.args.method, "http://localhost:8000", body)

			if tt.args.compressHeader != "" {
				req.Header.Add("content-encoding", tt.args.compressHeader)
			}

			res := httptest.NewRecorder()
			handler(res, req)

			got := GzipDecompressMiddleware(handler)
			got.ServeHTTP(res, req)
			result := res.Result()
			defer result.Body.Close()

			require.Equal(t, tt.wantStatus, result.StatusCode,
				fmt.Sprintf("wrong status code = %v, want %v", result.StatusCode, tt.wantStatus))

		})
	}
}

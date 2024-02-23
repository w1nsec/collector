package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w1nsec/collector/internal/config"
)

func TestCIDRmiddleware_Handle(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	_, validCidr, _ := net.ParseCIDR("10.0.0.0/24")

	type args struct {
		method     string
		realIP     []string
		respStatus int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test Valid cidr",
			args: args{
				realIP:     []string{"10.0.0.125", "11.12.3.14"},
				respStatus: http.StatusOK,
			},
		},
		{
			name: "Test Valid cidr2",
			args: args{
				realIP:     []string{"11.12.3.14", "10.0.0.125"},
				respStatus: http.StatusOK,
			},
		},
		{
			name: "Test Valid cidr POST",
			args: args{
				method:     http.MethodPost,
				realIP:     []string{"10.0.0.125", "11.12.3.14"},
				respStatus: http.StatusOK,
			},
		},
		{
			name: "Test Valid cidr2 POST",
			args: args{
				method:     http.MethodPost,
				realIP:     []string{"11.12.3.14", "10.0.0.125"},
				respStatus: http.StatusOK,
			},
		},
		{
			name: "Test Invalid cidr",
			args: args{
				realIP:     []string{"11.12.3.14", "10.1.11.125"},
				respStatus: http.StatusForbidden,
			},
		},
		{
			name: "Test Invalid cidr POST",
			args: args{
				method:     http.MethodPost,
				realIP:     []string{"11.12.3.14", "10.1.11.125"},
				respStatus: http.StatusForbidden,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &CIDRmiddleware{
				cidr: validCidr,
			}

			req := httptest.NewRequest(tt.args.method, "http://localhost:18239", nil)
			//req.Header.Add(config.RealIPHeader, strings.Join(tt.args.realIP, "; ")+";")
			//req.Header[config.RealIPHeader] = tt.args.realIP
			for _, ip := range tt.args.realIP {
				req.Header.Add(config.RealIPHeader, ip)
			}

			recorder := httptest.NewRecorder()
			handler(recorder, req)

			h := m.Handle(handler)
			h.ServeHTTP(recorder, req)

			res := recorder.Result()
			defer res.Body.Close()

			require.Equal(t, tt.args.respStatus, res.StatusCode,
				fmt.Sprintf("wrong status code = %v, want %v",
					res.StatusCode, tt.args.respStatus))

		})
	}
}

func TestNewCIDRmiddleware(t *testing.T) {
	var addr = "192.168.1.1/24"
	_, cidr, _ := net.ParseCIDR(addr)
	middleware := NewCIDRmiddleware(cidr)

	require.NotNil(t, middleware, "return nil from constructor")
	require.Equal(t, *cidr, *middleware.cidr)
}

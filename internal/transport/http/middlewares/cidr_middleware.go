package middlewares

import (
	"net"
	"net/http"

	"github.com/w1nsec/collector/internal/config"
)

type CIDRmiddleware struct {
	cidr *net.IPNet
}

func NewCIDRmiddleware(cidr *net.IPNet) *CIDRmiddleware {
	return &CIDRmiddleware{cidr: cidr}
}

func (m *CIDRmiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ips := r.Header.Values(config.RealIPHeader)
		var contains = false
		for _, ipStr := range ips {
			ip := net.ParseIP(ipStr)
			if m.cidr.Contains(ip) {
				contains = true
			}
		}

		if !contains {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

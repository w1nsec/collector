package middlewares

import "net/http"

func DebugStore(h http.Handler) http.Handler {
	debugFn := func(w http.ResponseWriter, r *http.Request) {

	}
	return http.HandlerFunc(debugFn)
}
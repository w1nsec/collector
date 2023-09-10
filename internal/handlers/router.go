package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/w1nsec/collector/internal/memstorage"
	"net/http"
)

func NewRouter(store memstorage.Storage) http.Handler {
	r := chi.NewRouter()
	r.Route("/update/", func(r chi.Router) {
		//r.Post("/", UpdateMetricsHandle(store))
		r.Post("/counter/{name}/{value}", UpdateCounterHandle(store))
		r.Post("/gauge/{name}/{value}", UpdateGaugeHandle(store))

		// Not Found
		r.Post("/gauge/{name}/", NotFoundHandle)
		r.Post("/counter/{name}/", NotFoundHandle)

		// Bad Request (other paths)
		// TODO CHANGE other paths to BadRequest
		// Now only /wrongtype/ return 405
		// Want: /wrongtype/metricname/123 return 405
		//       /wrongtype/ 			   return 405
		r.Post("/{other:[a-zA-Z0-9/]+}", BadRequest)
	})

	r.Route("/value/", func(r chi.Router) {
		r.Get("/{mType}/{mName}", GetMetric(store))
	})
	return r
}

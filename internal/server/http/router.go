package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/w1nsec/collector/internal/handlers"
	"github.com/w1nsec/collector/internal/middlewares"
	"github.com/w1nsec/collector/internal/service"
	"net/http"
)

func NewRouter(service service.Service) http.Handler {
	r := chi.NewRouter()

	// middlewares
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)

	// routing
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAllMetrics(service))
	})

	r.Route("/update/", func(r chi.Router) {
		//r.Use(printMidl)
		//r.Post("/", UpdateMetricsHandle(store))
		r.Post("/counter/{name}/{value}", handlers.UpdateCounterHandle(service))
		r.Post("/gauge/{name}/{value}", handlers.UpdateGaugeHandle(service))

		// Update via JSON (only one metric by yandex TASK)
		//r.Post("/", JSONUpdateHandler(store))
		r.Post("/", handlers.JSONUpdateOneMetricHandler(service))

		// Not Found
		r.Post("/gauge/", handlers.NotFoundHandle)
		r.Post("/counter/", handlers.NotFoundHandle)
		r.Post("/gauge/{name}", handlers.NotFoundHandle)
		r.Post("/counter/{name}", handlers.NotFoundHandle)
		//r.Post("/gauge/{name}/", NotFoundHandle)
		//r.Post("/counter/{name}/", NotFoundHandle)

		// Bad Request (other paths)
		// TODO CHANGE other paths to BadRequest
		// Now only /wrongtype/ return 405
		// Want: /wrongtype/metricname/123 return 405
		//       /wrongtype/ 			   return 405
		//r.Post("/{other}", BadRequest)
		r.NotFound(handlers.BadRequest)
	})

	r.Route("/value/", func(r chi.Router) {
		//r.Use(middlewares.LoggingMiddleware)

		// Get metric value
		r.Post("/", handlers.JSONGetMetricHandler(service))
		r.Get("/{mType}/{mName}", handlers.GetMetric(service))
	})

	/// increment 6 testing
	r.Route("/echoping", func(r chi.Router) {
		//r.Use(middlewares.LoggingMiddleware)
		r.Get("/", handlers.Pong)
	})

	/// increment 10
	r.Route("/ping", func(r chi.Router) {
		//r.Use(middlewares.LoggingMiddleware)
		r.Get("/", handlers.CheckDBConnectionHandler(service))
	})

	return r
}

func printMidl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

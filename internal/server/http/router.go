package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	handlers2 "github.com/w1nsec/collector/internal/server/http/handlers"
	"github.com/w1nsec/collector/internal/server/http/middlewares"
	"github.com/w1nsec/collector/internal/service"
	"net/http"
)

func NewRouter(service *service.MetricService) http.Handler {
	r := chi.NewRouter()

	// middlewares
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)

	// routing
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers2.GetAllMetrics(service))
	})

	r.Route("/update/", func(r chi.Router) {
		//r.Use(printMidl)
		//r.Post("/", UpdateMetricsHandle(store))
		r.Post("/counter/{name}/{value}", handlers2.UpdateCounterHandle(service))
		r.Post("/gauge/{name}/{value}", handlers2.UpdateGaugeHandle(service))

		// Update via JSON (only one metric by yandex TASK)
		//r.Post("/", JSONUpdateHandler(store))
		r.Post("/", handlers2.JSONUpdateOneMetricHandler(service))

		// Not Found
		r.Post("/gauge/", handlers2.NotFoundHandle)
		r.Post("/counter/", handlers2.NotFoundHandle)
		r.Post("/gauge/{name}", handlers2.NotFoundHandle)
		r.Post("/counter/{name}", handlers2.NotFoundHandle)
		//r.Post("/gauge/{name}/", NotFoundHandle)
		//r.Post("/counter/{name}/", NotFoundHandle)

		// Bad Request (other paths)
		// TODO CHANGE other paths to BadRequest
		// Now only /wrongtype/ return 405
		// Want: /wrongtype/metricname/123 return 405
		//       /wrongtype/ 			   return 405
		//r.Post("/{other}", BadRequest)
		r.NotFound(handlers2.BadRequest)
	})

	r.Route("/value/", func(r chi.Router) {
		//r.Use(middlewares.LoggingMiddleware)

		// Get metric value
		r.Post("/", handlers2.JSONGetMetricHandler(service))
		r.Get("/{mType}/{mName}", handlers2.GetMetric(service))
	})

	/// increment 6 testing
	r.Route("/echoping", func(r chi.Router) {
		//r.Use(middlewares.LoggingMiddleware)
		r.Get("/", handlers2.Pong)
	})

	/// increment 10
	r.Route("/ping", func(r chi.Router) {
		//r.Use(middlewares.LoggingMiddleware)
		r.Get("/", handlers2.CheckDBConnectionHandler(service))
	})

	return r
}

func printMidl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

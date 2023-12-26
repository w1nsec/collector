package http

import (
	"github.com/go-chi/chi/v5"
	chimidl "github.com/go-chi/chi/v5/middleware"
	"github.com/w1nsec/collector/internal/service"
	"github.com/w1nsec/collector/internal/transport/http/handlers"
	"github.com/w1nsec/collector/internal/transport/http/middlewares"
	"net/http"
)

var defaultCompressibleContentTypes = []string{
	"text/html",
	"text/plain",
	"application/json",
}

func NewRouter(service *service.MetricService) http.Handler {
	r := chi.NewRouter()

	// middlewares

	signMidl := middlewares.NewSigningMidl(service.Secret)
	r.Use(signMidl.Signing)
	r.Use(middlewares.LoggingMiddleware)
	//r.Use(middlewares.GzipMiddleware)
	r.Use(middlewares.GzipDecompressMiddleware)
	r.Use(chimidl.Compress(5, defaultCompressibleContentTypes...))

	// handlers
	getAllMetrics := handlers.NewGetMetricsHandler(service)
	getOneMetric := handlers.NewGetMetricHandler(service)
	updateCounters := handlers.NewUpdateCountersHandler(service)
	updateGauges := handlers.NewUpdateGaugeHandler(service)
	jsonUpdateOne := handlers.NewJSONUpdateOneMetricHandler(service)
	jsonUpdateAll := handlers.NewJSONUpdateMetricsHandler(service)
	jsonGetMetric := handlers.NewJSONGetMetricHandler(service)
	dbCheck := handlers.NewCheckDBConnectionHandler(service)

	// routing
	r.Route("/", func(r chi.Router) {
		r.Get("/", getAllMetrics.ServeHTTP)
	})

	r.Route("/update/", func(r chi.Router) {
		r.Post("/counter/{name}/{value}", updateCounters.ServeHTTP)
		r.Post("/gauge/{name}/{value}", updateGauges.ServeHTTP)

		// Update via JSON (only one metric by yandex TASK)
		r.Post("/", jsonUpdateOne.ServeHTTP)

		// Not Found
		r.Group(func(r chi.Router) {
			r.Post("/gauge/", handlers.NotFoundHandle)
			r.Post("/gauge/", handlers.NotFoundHandle)
			r.Post("/counter/", handlers.NotFoundHandle)
			r.Post("/gauge/{name}", handlers.NotFoundHandle)
			r.Post("/counter/{name}", handlers.NotFoundHandle)
		})

		// Bad Request (other paths)
		// TODO CHANGE other paths to BadRequest
		// Now only /wrongtype/ return 405
		// Want: /wrongtype/metricname/123 return 405
		//       /wrongtype/ 			   return 405
		//r.Post("/{other}", BadRequest)
		r.NotFound(handlers.BadRequest)
	})

	r.Post("/updates/", jsonUpdateAll.ServeHTTP)

	r.Route("/value/", func(r chi.Router) {
		// Get metric value
		r.Post("/", jsonGetMetric.ServeHTTP)
		r.Get("/{mType}/{mName}", getOneMetric.ServeHTTP)
	})

	/// increment 6 testing
	r.Route("/echoping", func(r chi.Router) {
		r.Get("/", handlers.Pong)
	})

	/// increment 10
	r.Route("/ping", func(r chi.Router) {
		//r.Use(middlewares.LoggingMiddleware)
		r.Get("/", dbCheck.ServeHTTP)
	})

	// increment 16
	r.Mount("/debug", chimidl.Profiler())
	return r
}

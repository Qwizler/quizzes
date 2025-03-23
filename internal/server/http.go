package server

import (
	"github.com/gorilla/handlers"
	quizzesV1 "quiz/api/quizzes/v1"
	"quiz/internal/conf"
	"quiz/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func NewHTTPServer(
	c *conf.Server,
	quizzes *service.QuizzesService,
	logger log.Logger,
	meter metric.Meter,
	tp trace.TracerProvider,
) (*http.Server, error) {
	counter, err := metrics.DefaultRequestsCounter(meter, metrics.DefaultServerRequestsCounterName)
	if err != nil {
		return nil, err
	}
	seconds, err := metrics.DefaultSecondsHistogram(meter, metrics.DefaultServerSecondsHistogramName)
	if err != nil {
		return nil, err
	}
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(
				tracing.WithTracerProvider(tp),
			),
			logging.Server(logger),
			metrics.Server(
				metrics.WithRequests(counter),
				metrics.WithSeconds(seconds),
			),
		),
	}
	if c.Http.GetCors().GetEnabled() {
		allowHeaders := c.Http.GetCors().GetAllowHeaders()
		allowMethods := c.Http.GetCors().GetAllowMethods()
		allowOrigins := c.Http.GetCors().GetAllowOrigins()
		cors := handlers.CORS(
			handlers.AllowedHeaders(allowHeaders),
			handlers.AllowedMethods(allowMethods),
			handlers.AllowedOrigins(allowOrigins),
		)
		opts = append(opts, http.Filter(cors))
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	srv.HandlePrefix("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))

	quizzesV1.RegisterQuizzesHTTPServer(srv, quizzes)
	return srv, nil
}

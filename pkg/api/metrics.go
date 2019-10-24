package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	reqMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "helloworld_http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"endpoint", "code"})

	reqDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "helloworld_http_request_duration_seconds",
		Help: "Duration of the HTTP requests",
	}, []string{"endpoint"})
)

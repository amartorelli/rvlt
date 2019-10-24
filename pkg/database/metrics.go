package database

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	opMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "helloworld_database_operations_total",
		Help: "Total number of database operations",
	}, []string{"operation"})

	opErrMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "helloworld_database_operations_errors_total",
		Help: "Total number of database operation errors",
	}, []string{"operation"})

	opDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "helloworld_database_operations_duration_seconds",
		Help: "Duration of the database operations",
	}, []string{"operation"})
)

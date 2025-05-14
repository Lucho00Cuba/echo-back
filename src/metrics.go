package main

import "github.com/prometheus/client_golang/prometheus"

var (
	// requestCount is a Prometheus counter for the number of requests served.
	requestCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "request_count",
		Help: "Number of requests served.",
	})
	// requestDuration is a Prometheus histogram for request duration in seconds.
	requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Duration of requests in seconds.",
		Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
	})
	// requestStatusCount is a Prometheus counter for the number of requests served by status code.
	requestStatusCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status_total",
			Help: "Count of responses by HTTP status code",
		},
		[]string{"code"},
	)
	// requestDurationVec is a Prometheus histogram for request duration in seconds by status code.
	requestDurationVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests by status code",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "code"},
	)
)

// init initializes the Prometheus metrics.
func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestStatusCount)
	prometheus.MustRegister(requestDurationVec)
}

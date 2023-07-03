package exporter

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	Size     prometheus.Counter
	Duration *prometheus.HistogramVec
	Requests *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		Size: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "nginx",
			Name:      "size_bytes_total",
			Help:      "Total bytes sent to the clients.",
		}),
		Requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "nginx",
			Name:      "http_requests_total",
			Help:      "Total number of requests.",
		}, []string{"status_code", "method", "path"}),
		Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "nginx",
			Name:      "http_request_duration_seconds",
			Help:      "Duration of the request.",
			// Optionally configure time buckets
			// Buckets:   prometheus.LinearBuckets(0.01, 0.05, 20),
			Buckets: prometheus.DefBuckets,
		}, []string{"status_code", "method", "path"}),
	}
	reg.MustRegister(m.Size, m.Requests, m.Duration)
	return m
}

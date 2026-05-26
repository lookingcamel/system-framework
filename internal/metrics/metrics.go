package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	RequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)

	AppInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_info",
			Help: "Application information",
		},
		[]string{"name", "version"},
	)

	HealthStatus = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "health_status",
			Help: "Health check status (1 = healthy, 0 = unhealthy)",
		},
	)
)

func RecordRequest(method, path, status string, duration float64) {
	RequestsTotal.WithLabelValues(method, path, status).Inc()
	RequestDuration.WithLabelValues(method, path).Observe(duration)
}

func IncrementInFlight() {
	RequestsInFlight.Inc()
}

func DecrementInFlight() {
	RequestsInFlight.Dec()
}

func SetAppInfo(name, version string) {
	AppInfo.WithLabelValues(name, version).Set(1)
}

func SetHealthStatus(healthy bool) {
	if healthy {
		HealthStatus.Set(1)
	} else {
		HealthStatus.Set(0)
	}
}

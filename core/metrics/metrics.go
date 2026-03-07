// Package metrics provides Prometheus HTTP instrumentation for Gin.
package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds Prometheus collectors for HTTP instrumentation.
type Metrics struct {
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
	size     *prometheus.HistogramVec
}

// New creates a Metrics instance and registers collectors with the default
// Prometheus registry.
func New() *Metrics {
	return newMetrics(prometheus.DefaultRegisterer)
}

// newMetrics creates a Metrics instance with the given registerer.
// Used by tests to provide an isolated prometheus.NewRegistry().
func newMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		}, []string{"method", "path", "status"}),

		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path", "status"}),

		size: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes.",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		}, []string{"method", "path", "status"}),
	}

	reg.MustRegister(m.requests, m.duration, m.size)
	return m
}

// Middleware returns a Gin middleware that records request count, duration,
// and response size for every HTTP request.
func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		path := c.FullPath()
		if path == "" {
			path = "unmatched"
		}
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		elapsed := time.Since(start).Seconds()

		m.requests.WithLabelValues(method, path, status).Inc()
		m.duration.WithLabelValues(method, path, status).Observe(elapsed)
		m.size.WithLabelValues(method, path, status).Observe(float64(c.Writer.Size()))
	}
}

// Handler returns a Gin handler that serves Prometheus metrics in the text
// exposition format.
func Handler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

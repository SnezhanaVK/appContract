package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func MonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path

		rw := &responseWriterWrapper{w: w}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		status := fmt.Sprintf("%d", rw.status)
		httpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
	})
}

type responseWriterWrapper struct {
	w      http.ResponseWriter
	status int
}

func (rw *responseWriterWrapper) Header() http.Header {
	return rw.w.Header()
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	return rw.w.Write(b)
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.w.WriteHeader(statusCode)
}

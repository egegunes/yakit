package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requests     *prometheus.CounterVec
	responseTime *prometheus.HistogramVec
)

func init() {
	responseTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "response_time_seconds",
		Help:      "Request response times",
	}, []string{"code", "method", "path"})
	prometheus.MustRegister(responseTime)

	requests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "http",
		Name:      "http_requests_total",
		Help:      "HTTP requests processed",
	}, []string{"code", "method", "path"})
	prometheus.MustRegister(requests)
}

type metricsWriter struct {
	w      http.ResponseWriter
	status int
}

func (m metricsWriter) Write(b []byte) (int, error) {
	return m.w.Write(b)
}

func (m metricsWriter) Header() http.Header {
	return m.w.Header()
}

func (m *metricsWriter) WriteHeader(statusCode int) {
	m.status = statusCode
	m.w.WriteHeader(statusCode)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		m := &metricsWriter{w: w}

		next.ServeHTTP(m, r)

		status := strconv.Itoa(m.status)

		requests.WithLabelValues(status, r.Method, r.URL.Path).Inc()
		responseTime.WithLabelValues(status, r.Method, r.URL.Path).Observe(float64(time.Since(start).Seconds()))
	})
}

package middlewares

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests.",
		},
		[]string{"method"},
	)
)

func Metrics(next http.Handler) http.Handler {
	prometheus.MustRegister(requestsTotal, requestDuration)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func() {
			method := r.Method
			elapsed := time.Since(start).Seconds()
			requestsTotal.WithLabelValues(method).Inc()
			requestDuration.WithLabelValues(method).Observe(elapsed)
		}()

		next.ServeHTTP(w, r)
	})
}

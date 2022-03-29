package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsMiddleware struct {
	opsProcessed *prometheus.CounterVec
	opsDuration  *prometheus.HistogramVec
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func NewMetricsMiddleware() *MetricsMiddleware {
	labels := []string{"method", "path", "status_code"}

	opsProcessed := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "myapp",
		Name:      "request_counter",
		Help:      "Counter of processed requests",
	}, labels)

	buckets := []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	opsDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "myapp",
		Name:      "request_latency_seconds",
		Help:      "Histogram of request processing time latency in seconds",
		Buckets:   buckets,
	}, labels)

	return &MetricsMiddleware{
		opsProcessed: opsProcessed,
		opsDuration:  opsDuration,
	}
}

func (mm *MetricsMiddleware) Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := statusRecorder{w, 200}

		next.ServeHTTP(&rec, r)

		duration := time.Since(start)
		statusCode := strconv.Itoa(rec.statusCode)
		path := r.URL.Path
		labels := prometheus.Labels{"method": r.Method, "path": path, "status_code": statusCode}

		mm.opsProcessed.With(labels).Inc()
		mm.opsDuration.With(labels).Observe(duration.Seconds())
	})
}

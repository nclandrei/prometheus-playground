package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type CustomError struct {
	Code    int
	Message string
	Detail  error
	Stack   string
}

var (
	histo = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "hello_world_duration_seconds",
		Help:    "Describes the latency of a HTTP handler.",
		Buckets: []float64{0.05, 0.1, 0.5, 1, 100},
	}, []string{"method", "code"})
	errMetrics = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hello_world_errors",
		Help: "Describes the errors retrieved through the HTTP handler.",
	}, []string{"code", "message", "detail", "stack"})
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		err := CustomError{
			Code:    10002,
			Message: "This is a test error 2",
			Detail:  fmt.Errorf("wrapper Error here, dude"),
			Stack:   "happened at line -1.",
		}
		duration := time.Since(start)
		code := 200
		errMetrics.WithLabelValues(fmt.Sprintf("%d", err.Code), err.Message,
			err.Detail.Error(), err.Stack).Inc()
		histo.WithLabelValues(r.Method, fmt.Sprintf("%d", code)).Observe(duration.Seconds())
	})
	http.Handle("/metrics", prometheus.Handler())
	prometheus.MustRegister(histo, errMetrics)

	http.ListenAndServe(":3000", nil)
}

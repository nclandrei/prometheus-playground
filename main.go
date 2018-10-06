package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	histo = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "hello_world_duration_seconds",
		Help: "Describes the latency of a HTTP handler.",
	}, []string{"method", "code"})
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Fprintf(w, "Hello world!")
		duration := time.Since(start)
		code := 200
		histo.WithLabelValues(r.Method, fmt.Sprintf("%d", code)).Observe(duration.Seconds())
	})
	http.Handle("/metrics", prometheus.Handler())
	prometheus.Register(histo)

	http.ListenAndServe(":3000", nil)
}

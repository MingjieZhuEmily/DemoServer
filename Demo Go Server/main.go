package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type opType int

const (
	fast opType = iota
	slow
	failed
)

func main() {
	// Serve the front-end files
	http.Handle("/", http.FileServer(http.Dir("static")))

	// Register Prometheus handler for metrics collection
	http.Handle("/metrics", promhttp.Handler())

	handler := newHandler()

	// Set up routes to handle API requests
	http.HandleFunc("/fast", handler.handleFastRequest)
	http.HandleFunc("/slow", handler.handleSlowRequest)
	http.HandleFunc("/failed", handler.handleFailedRequest)

	// Start the server
	fmt.Println("Server listening on localhost:8080")
	http.ListenAndServe(":8080", nil)
}

type handler struct {
	histogram *prometheus.HistogramVec
}

func newHandler() *handler {
	hist := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "demo_app_http_request_duration_seconds",
		Help: "Request duration in seconds",
	}, []string{"user", "code"})

	return &handler{histogram: hist}
}

func (h *handler) handleFastRequest(w http.ResponseWriter, r *http.Request) {
	h.performRequest(w, r, fast)
}

func (h *handler) handleSlowRequest(w http.ResponseWriter, r *http.Request) {
	h.performRequest(w, r, slow)
}

func (h *handler) handleFailedRequest(w http.ResponseWriter, r *http.Request) {
	h.performRequest(w, r, failed)
}

func (h *handler) performRequest(w http.ResponseWriter, r *http.Request, typ opType) {
	username := r.FormValue("username")

	switch typ {
	case fast:
		h.recordMetric(username, time.Duration(10)*time.Millisecond, 200)

	case slow:
		h.recordMetric(username, time.Duration(10)*time.Second, 200)

	case failed:
		h.recordMetric(username, time.Duration(10)*time.Millisecond, 500)
	}

	// Return a response to the client
	fmt.Fprintf(w, "Request completed successfully!")
}

func (h *handler) recordMetric(username string, duration time.Duration, code int) {
	h.histogram.WithLabelValues(username, fmt.Sprintf("%d", code)).Observe(duration.Seconds())
}

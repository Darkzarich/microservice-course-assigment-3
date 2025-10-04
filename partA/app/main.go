
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "code"},
	)
	reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "Request duration seconds",
			Buckets: []float64{0.05, 0.1, 0.2, 0.5, 1, 2, 5},
		},
		[]string{"path"},
	)
)

func record(w http.ResponseWriter, r *http.Request, status int, path string, start time.Time) {
	elapsed := time.Since(start).Seconds()
	reqDuration.WithLabelValues(path).Observe(elapsed)
	httpRequests.WithLabelValues(r.Method, path, fmt.Sprintf("%d", status)).Inc()
}

func handlerOK(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	// random latency
	// about 20% slowed ones (>0.5s)
	slowProb := 0.2
	if rand.Float64() < slowProb {
		time.Sleep(time.Duration(600+rand.Intn(600)) * time.Millisecond)
	} else {
		time.Sleep(time.Duration(10+rand.Intn(80)) * time.Millisecond)
	}

	msg := fmt.Sprintf(`{"status":"ok","time":"%s"}`, time.Now().Format(time.RFC3339))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(msg))
	record(w, r, http.StatusOK, "/", start)
}

func handlerError(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	// simulate error
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(`{"error":"demo failure"}`))
	record(w, r, http.StatusInternalServerError, "/error", start)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(reqDuration)

	http.HandleFunc("/", handlerOK)
	http.HandleFunc("/error", handlerError)

	// metrics endpoint
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		addr := ":19090"
		log.Printf("metrics on %s", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatalf("metrics server error: %v", err)
		}
	}()

	addr := ":8080"
	if v := os.Getenv("APP_PORT"); v != "" {
		addr = ":" + v
	}
	log.Printf("app listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("http server error: %v", err)
	}
}

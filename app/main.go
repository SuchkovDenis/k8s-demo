package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

// Max concurrent requests before returning 503
const maxConcurrent = 5

var inFlight atomic.Int32

func main() {
	hostname, _ := os.Hostname()
	fmt.Printf("Server starting on :8080 (pod: %s)\n", hostname)

	http.HandleFunc("/work", cors(workHandler))
	http.HandleFunc("/health", cors(healthHandler))
	http.ListenAndServe(":8080", nil)
}

func workHandler(w http.ResponseWriter, r *http.Request) {
	// Reject if overloaded
	current := inFlight.Add(1)
	defer inFlight.Add(-1)

	hostname, _ := os.Hostname()
	w.Header().Set("Content-Type", "application/json")

	if current > maxConcurrent {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"pod":   hostname,
			"error": "overloaded",
		})
		return
	}

	start := time.Now()

	// CPU-intensive work
	result := 0.0
	for i := 0; i < 50_000_000; i++ {
		result += math.Sqrt(float64(i))
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"pod":    hostname,
		"timeMs": time.Since(start).Milliseconds(),
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next(w, r)
	}
}

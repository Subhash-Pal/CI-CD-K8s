package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type quoteResponse struct {
	Quote     string    `json:"quote"`
	Source    string    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", healthHandler)
	mux.HandleFunc("/quote", quoteHandler)

	port := getenv("QUOTE_PORT", "8081")
	log.Printf("quote-service listening on: %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "quote-service",
	})
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, quoteResponse{
		Quote:     "Kubernetes makes reliable deployments repeatable.",
		Source:    "go-cd-demo",
		Timestamp: time.Now().UTC(),
	})
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("encode response: %v", err)
	}
}

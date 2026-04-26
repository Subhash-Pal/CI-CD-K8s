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

type apiResponse struct {
	Service   string        `json:"service"`
	Message   string        `json:"message"`
	Quote     quoteResponse `json:"quote"`
	Timestamp time.Time     `json:"timestamp"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", healthHandler)
	mux.HandleFunc("/api", apiHandler)

	port := getenv("API_PORT", "8080")
	log.Printf("api listening on %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "ok",
		"service": "api",
	}
	writeJSON(w, http.StatusOK, response)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	quoteURL := getenv("QUOTE_SERVICE_URL", "http://quote-service:8081/quote")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(quoteURL)
	if err != nil {
		http.Error(w, "quote service unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	var quote quoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		http.Error(w, "invalid quote payload", http.StatusBadGateway)
		return
	}

	response := apiResponse{
		Service:   "api",
		Message:   "CI/CD and Kubernetes demo is running",
		Quote:     quote,
		Timestamp: time.Now().UTC(),
	}
	writeJSON(w, http.StatusOK, response)
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

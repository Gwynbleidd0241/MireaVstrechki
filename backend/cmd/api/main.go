package main

import (
	"encoding/json"
	"log"
	"meeting-service/backend/internal/config"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)

	addr := ":" + cfg.Server.Port
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := healthResponse{
		Status: "ok",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

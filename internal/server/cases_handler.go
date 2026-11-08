package server

import (
	"detective/internal/cases"
	"encoding/json"
	"log"
	"net/http"
)

func casesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	registry := cases.NewCaseRegistry()

	availableCases := registry.List()

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(availableCases); err != nil {
		log.Printf("cases hadler: encode response: %v", err)
	}
}

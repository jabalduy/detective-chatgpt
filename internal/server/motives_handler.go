package server

import (
	"detective/internal/game"
	"encoding/json"
	"log"
	"net/http"
)

func motivesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, err := getGame(w, r)
	if err != nil {
		log.Printf("motives: get game: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	session.Mu.RLock()
	defer session.Mu.RUnlock()

	state := session.State
	if state == nil {
		http.Error(w, "game not started", http.StatusConflict)
		return
	}

	foundMotives := make([]game.Motive, 0)

	for _, motive := range state.Motives {
		if state.FoundMotives[motive.ID] {
			foundMotives = append(foundMotives, motive)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(foundMotives); err != nil {
		log.Printf("motives: encode response: %v", err)
	}
}

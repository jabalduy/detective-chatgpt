package server

import (
	"detective/internal/game"
	"encoding/json"
	"log"
	"net/http"
)

func factsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, err := getGame(w, r)
	if err != nil {
		log.Printf("facts: get game: %v", err)
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

	foundFacts := make([]game.Fact, 0)

	for _, fact := range state.Facts {
		if state.FoundFacts[fact.ID] {
			foundFacts = append(foundFacts, fact)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(foundFacts); err != nil {
		log.Printf("facts: encode response: %v", err)
	}
}

package server

import (
	"detective/internal/game"
	"encoding/json"
	"log"
	"net/http"
)

type DeductionRequest struct {
	FactIDs []int `json:"fact_ids"`
}

type DeductionResponse struct {
	Success    bool          `json:"success"`
	Message    string        `json:"message"`
	NewFacts   []game.Fact   `json:"new_facts,omitempty"`
	NewMotives []game.Motive `json:"new_motives,omitempty"`
}

func deductionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request DeductionRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(request.FactIDs) < 2 {
		http.Error(w, "at least two facts are required", http.StatusBadRequest)
		return
	}

	session, err := getGame(w, r)
	if err != nil {
		log.Printf("deduction: get game: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	session.Mu.Lock()
	defer session.Mu.Unlock()

	state := session.State

	if state == nil {
		http.Error(w, "game not started", http.StatusConflict)
		return
	}

	// 1. Ищем дедукцию по выбранным фактам
	deduction := game.FindDeduction(state, request.FactIDs)

	if deduction == nil {
		writeDeductionResponse(w, DeductionResponse{
			Success: false,
			Message: "Эти факты не приводят к новому выводу",
		})
		return
	}

	// 2. ЗАПОМИНАЕМ состояние ДО дедукции

	factsBefore := make(map[int]bool)

	for id, found := range state.FoundFacts {
		if found {
			factsBefore[id] = true
		}
	}

	motivesBefore := make(map[int]bool)

	for id, found := range state.FoundMotives {
		if found {
			motivesBefore[id] = true
		}
	}

	// 3. ВЫПОЛНЯЕМ дедукцию
	if !game.SolveDeduction(state, *deduction) {
		writeDeductionResponse(w, DeductionResponse{
			Success: false,
			Message: "Эту дедукцию сейчас невозможно выполнить",
		})
		return
	}

	// 4. СМОТРИМ, какие факты появились ПОСЛЕ дедукции
	newFacts := make([]game.Fact, 0)

	for _, fact := range state.Facts {
		if state.FoundFacts[fact.ID] && !factsBefore[fact.ID] {
			newFacts = append(newFacts, fact)
		}
	}

	// 5. СМОТРИМ, какие мотивы появились ПОСЛЕ дедукции
	newMotives := make([]game.Motive, 0)

	for _, motive := range state.Motives {
		if state.FoundMotives[motive.ID] && !motivesBefore[motive.ID] {
			newMotives = append(newMotives, motive)
		}
	}

	// 6. Возвращаем результат
	writeDeductionResponse(w, DeductionResponse{
		Success:    true,
		Message:    "Новый вывод получен",
		NewFacts:   newFacts,
		NewMotives: newMotives,
	})
}

func writeDeductionResponse(w http.ResponseWriter, response DeductionResponse) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("deduction: encode response: %v", err)
	}
}

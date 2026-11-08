package server

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
)

type AccusationRequest struct {
	SuspectID int   `json:"suspect_id"`
	MotiveID  int   `json:"motive_id"`
	FactIDs   []int `json:"fact_ids"`
}

// ОБВИНЕНИЕ
func accuseHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getGame(w, r)
	if err != nil {
		log.Printf("new game: %v", err)
		http.Error(w, "Не удалось загрузить игру", http.StatusInternalServerError)
		return
	}

	session.Mu.RLock()
	defer session.Mu.RUnlock()
	state := session.State

	if state == nil {
		http.Error(w, "Игра не запущена", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Некорректная команда", http.StatusMethodNotAllowed)
		return
	}

	var accusationRequest AccusationRequest
	err = json.NewDecoder(r.Body).Decode(&accusationRequest)
	if err != nil {
		http.Error(
			w,
			"Не удалось прочитать accusationRequest",
			http.StatusBadRequest,
		)
		return
	}

	if !state.FoundMotives[accusationRequest.MotiveID] {
		http.Error(
			w,
			"Отправленный мотив еще не был найден",
			http.StatusBadRequest,
		)
		return
	}

	for _, factID := range accusationRequest.FactIDs {
		if !state.FoundFacts[factID] {
			http.Error(
				w,
				"Один из отправленных фактов еще не был найден",
				http.StatusBadRequest,
			)
			return
		}
	}

	foundMotive := false
	for _, motive := range state.Motives {
		if motive.ID == accusationRequest.MotiveID {
			foundMotive = true

			if motive.SuspectID != accusationRequest.SuspectID {
				http.Error(
					w,
					"Мотив принадлежит другому персонажу",
					http.StatusBadRequest,
				)
				return
			}

			break
		}
	}

	if !foundMotive {
		http.Error(
			w,
			"Мотив не найден",
			http.StatusBadRequest,
		)
		return
	}

	solution := state.Case.Solution

	allFactsCorrect := true
	for _, rightFact := range solution.FactIDs {
		if !slices.Contains(accusationRequest.FactIDs, rightFact) {
			allFactsCorrect = false
			break
		}
	}

	response := AccuseResponse{}
	switch {
	case accusationRequest.SuspectID != solution.SuspectID:
		response.TypeEnd = "fail"
		response.TextEnd = state.Case.Endings.Fail
	case accusationRequest.MotiveID != solution.MotiveID:
		response.TypeEnd = "unsolved"
		response.TextEnd = state.Case.Endings.Unsolved
	case !allFactsCorrect:
		response.TypeEnd = "unsolved"
		response.TextEnd = state.Case.Endings.Unsolved
	default:
		response.TypeEnd = "win"
		response.TextEnd = state.Case.Endings.Win
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Не удалось отправить обвинение", http.StatusInternalServerError)
		return
	}
}

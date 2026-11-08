package server

import (
	"detective/internal/game"
	"encoding/json"
	"log"
	"net/http"
)

// ДОСЬЕ
func caseHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getGame(w, r)
	if err != nil {
		log.Printf("new state: %v", err)
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

	response := CaseResponse{}

	for _, clue := range state.Clues {
		if state.FoundClues[clue.ObjID] {
			response.FoundClues = append(response.FoundClues, clue)
			response.CluesFound++
		}
	}

	for _, dial := range state.Dialogues {
		dialKey := game.DialogueID{
			SusID:  dial.SusID,
			DialID: dial.DialID,
		}

		if state.AskedDialogues[dialKey] {
			response.KnownFacts = append(response.KnownFacts, dial)
			response.DialoguesAsked++
		}
	}

	for _, obj := range state.Objects {
		if state.SearchedObjects[obj.ObjID] {
			response.ObjectsSearched++
		}
	}

	response.CluesTotal = len(state.Clues)
	response.DialoguesTotal = len(state.Dialogues)
	response.ObjectsTotal = len(state.Objects)

	response.Intro = state.Case.Intro
	response.KnownInfo = state.Case.KnownInfo

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Не удалось отправить досье", http.StatusInternalServerError)
		return
	}
}

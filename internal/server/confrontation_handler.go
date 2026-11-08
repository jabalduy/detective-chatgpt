package server

import (
	"detective/internal/game"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type ConfrontationResponse struct {
	Success bool   `json:"success"`
	Answer  string `json:"answer"`
}

func confrontationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Некорректная команда", http.StatusMethodNotAllowed)
		return
	}
	session, err := getGame(w, r)
	if err != nil {
		log.Printf("confrontation: get game: %v", err)
		http.Error(w, "Не удалось загрузить игру", http.StatusInternalServerError)
		return
	}
	session.Mu.Lock()
	defer session.Mu.Unlock()
	state := session.State
	if state == nil {
		http.Error(w, "Игра не запущена", http.StatusConflict)
		return
	}
	susID, err := strconv.Atoi(r.URL.Query().Get("susID"))
	if err != nil {
		http.Error(w, "Некорректный susID", http.StatusBadRequest)
		return
	}
	dialID, err := strconv.Atoi(r.URL.Query().Get("dialID"))
	if err != nil {
		http.Error(w, "Некорректный dialID", http.StatusBadRequest)
		return
	}
	for _, dialogue := range state.Dialogues {
		if dialogue.SusID != susID || dialogue.DialID != dialID {
			continue
		}
		if !game.PerformConfrontation(state, dialogue) {
			http.Error(w, "Конфронтация недоступна", http.StatusConflict)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ConfrontationResponse{Success: true, Answer: dialogue.Confrontation.Answer})
		return
	}
	http.Error(w, "Диалог не найден", http.StatusNotFound)
}

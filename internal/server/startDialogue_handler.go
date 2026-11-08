package server

import (
	"detective/internal/game"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func startDialogueHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getGame(w, r)
	if err != nil {
		log.Printf("new game: %v", err)
		http.Error(w, "Не удалось загрузить игру", http.StatusInternalServerError)
		return
	}

	session.Mu.Lock()
	defer session.Mu.Unlock()
	state := session.State

	if state == nil {
		http.Error(w, "Игра не запущена", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Некорректная команда", http.StatusMethodNotAllowed)
		return
	}

	susIDstr := r.URL.Query().Get("susID")
	susID, err := strconv.Atoi(susIDstr)
	if err != nil {
		http.Error(w, "Некорректный susID", http.StatusBadRequest)
		return
	}

	dialIDstr := r.URL.Query().Get("dialID")
	dialID, err := strconv.Atoi(dialIDstr)
	if err != nil {
		http.Error(w, "Некорректный dialID", http.StatusBadRequest)
		return
	}

	var response game.Dialogue

	found := false
	for _, dialogue := range state.Dialogues {
		if dialogue.SusID == susID &&
			dialogue.DialID == dialID {

			if !state.OpenDialogues[dialogue.ID()] {
				http.Error(w, "Диалог не доступен", http.StatusForbidden)
				return
			}

			if state.AskedDialogues[dialogue.ID()] {
				http.Error(w, "Диалог уже был", http.StatusBadRequest)
				return
			}

			if state.ActiveDialogue != nil {
				http.Error(w, "Диалог уже открыт", http.StatusConflict)
				return
			}

			id := dialogue.ID()
			state.ActiveDialogue = &id

			response = dialogue
			found = true

			state.Clock.Pause()

			break
		}
	}

	if !found {
		http.Error(w, "Не найден диалог", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Не удалось отправить выбранный диалог", http.StatusInternalServerError)
		return
	}
}

func finishDialogueHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getGame(w, r)
	if err != nil {
		log.Printf("new game: %v", err)
		http.Error(w, "Не удалось загрузить игру", http.StatusInternalServerError)
		return
	}

	session.Mu.Lock()
	defer session.Mu.Unlock()
	state := session.State

	if state == nil {
		http.Error(w, "Игра не запущена", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Некорректная команда", http.StatusMethodNotAllowed)
		return
	}

	susIDstr := r.URL.Query().Get("susID")
	susID, err := strconv.Atoi(susIDstr)
	if err != nil {
		http.Error(w, "Некорректный susID", http.StatusBadRequest)
		return
	}

	dialIDstr := r.URL.Query().Get("dialID")
	dialID, err := strconv.Atoi(dialIDstr)
	if err != nil {
		http.Error(w, "Некорректный dialID", http.StatusBadRequest)
		return
	}

	var response game.Dialogue

	found := false
	for _, dialogue := range state.Dialogues {
		if dialogue.SusID == susID &&
			dialogue.DialID == dialID {

			if !state.OpenDialogues[dialogue.ID()] {
				http.Error(w, "Диалог не доступен", http.StatusForbidden)
				return
			}

			if state.AskedDialogues[dialogue.ID()] {
				http.Error(w, "Диалог уже был", http.StatusBadRequest)
				return
			}

			if state.ActiveDialogue == nil {
				http.Error(w, "Диалог не запущен", http.StatusBadRequest)
				return
			}

			if *state.ActiveDialogue != dialogue.ID() {
				http.Error(w, "Попытка завершить другой диалог", http.StatusBadRequest)
				return
			}

			valid := game.PerformAction(state, dialogue.Action)
			if !valid {
				http.Error(w, "Не удалось применить действие", http.StatusInternalServerError)
				return
			}

			state.AskedDialogues[dialogue.ID()] = true
			state.ActiveDialogue = nil

			state.Clock.Resume()

			response = dialogue
			found = true

			break
		}
	}

	if !found {
		http.Error(w, "Не найден диалог", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Не удалось отправить выбранный диалог", http.StatusInternalServerError)
		return
	}
}

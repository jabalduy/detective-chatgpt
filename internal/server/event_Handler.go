package server

import (
	"detective/internal/game"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type EventResponse struct {
	ID          int                   `json:"event_id"`
	Title       string                `json:"event_title"`
	Description string                `json:"event_description"`
	TimeLimit   int                   `json:"time_limit"`
	Options     []EventOptionResponse `json:"event_options"`
}

type EventOptionResponse struct {
	ID   int    `json:"event_option_id"`
	Text string `json:"event_option_text"`
}

type ChooseEventResponse struct {
	Success bool `json:"success"`
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
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

	if r.Method != http.MethodGet {
		http.Error(w, "Некорректная команда", http.StatusMethodNotAllowed)
		return
	}

	if state.ActiveEvent == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	eventID := *state.ActiveEvent

	for _, event := range state.Events {
		if event.ID != eventID {
			continue
		}

		eventResponse := EventResponse{
			ID:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			TimeLimit:   event.TimeLimit,
			Options:     []EventOptionResponse{},
		}

		for _, option := range event.Options {
			eventOptionResponse := EventOptionResponse{
				ID:   option.ID,
				Text: option.Text,
			}

			eventResponse.Options = append(
				eventResponse.Options,
				eventOptionResponse,
			)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(eventResponse)
		if err != nil {
			http.Error(w, "Не удалось отправить ивент", http.StatusInternalServerError)
		}

		return
	}

	http.Error(
		w,
		"Активный ивент не найден",
		http.StatusInternalServerError,
	)
}

func chooseOptionHandler(w http.ResponseWriter, r *http.Request) {
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

	if state.ActiveEvent == nil {
		http.Error(w, "Нет активных ивентов", http.StatusConflict)
		return
	}

	optionIDstr := r.URL.Query().Get("optionID")

	optionID, err := strconv.Atoi(optionIDstr)
	if err != nil {
		http.Error(w, "Некорректный optionID", http.StatusBadRequest)
		return
	}

	if !game.ChooseEventOption(state, optionID) {
		http.Error(
			w,
			"Невозможно выбрать этот вариант",
			http.StatusConflict,
		)
		return
	}

	response := ChooseEventResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Не удалось отправить ответ", http.StatusInternalServerError)
		return
	}
}

func expireEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Некорректная команда", http.StatusMethodNotAllowed)
		return
	}
	session, err := getGame(w, r)
	if err != nil {
		http.Error(w, "Не удалось загрузить игру", http.StatusInternalServerError)
		return
	}
	session.Mu.Lock()
	defer session.Mu.Unlock()
	if session.State == nil {
		http.Error(w, "Игра не запущена", http.StatusBadRequest)
		return
	}
	session.State.ActiveEvent = nil
	w.WriteHeader(http.StatusNoContent)
}

package server

import (
	"detective/internal/cases"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// СТАРТ
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func startHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Некорректная команда", http.StatusMethodNotAllowed)
		return
	}

	caseID := r.URL.Query().Get("caseID")
	if caseID == "" {
		http.Error(w, "caseID is required", http.StatusBadRequest)
		return
	}

	state, err := newGame(caseID)
	if err != nil {
		if errors.Is(err, cases.ErrCaseNotFound) {
			log.Printf(
				"start game: unknown case: caseID=%q err=%v",
				caseID,
				err,
			)

			http.Error(
				w,
				"case not found",
				http.StatusNotFound,
			)
			return
		}

		log.Printf(
			"start game: caseID=%q err=%v",
			caseID,
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	session, err := getGame(w, r)
	if err != nil {
		log.Printf("start game: get session: %v", err)

		http.Error(
			w,
			"innternal server error",
			http.StatusInternalServerError,
		)
		return
	}

	session.Mu.Lock()
	session.State = state
	session.Mu.Unlock()

	response := StartResponse{
		Status:    "started",
		Title:     state.Case.Title,
		Message:   "Расследование начато",
		Intro:     state.Case.Intro,
		KnownInfo: state.Case.KnownInfo,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("start game: encode response: %v", err)
	}
}

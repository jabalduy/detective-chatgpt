package server

import (
	"crypto/rand"
	"detective/internal/cases"
	"detective/internal/game"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

var games = make(map[string]*Session)
var gamesMu sync.Mutex

var ErrGameNotStarted = errors.New("game not started")

type Session struct {
	State *game.GameState
	Mu    sync.RWMutex
}

func newSessionID() string {
	bytes := make([]byte, 16)

	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(bytes)
}

func getSessionID(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("sessionID")

	if err == nil {
		return cookie.Value
	}

	sessionID := newSessionID()

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionID",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
	})

	return sessionID
}

func getGame(w http.ResponseWriter, r *http.Request) (*Session, error) {
	sessionID := getSessionID(w, r)

	gamesMu.Lock()
	defer gamesMu.Unlock()

	session, ok := games[sessionID]
	if ok {
		return session, nil
	}

	session = &Session{
		State: nil,
	}

	games[sessionID] = session

	return session, nil
}

func newGame(caseID string) (*game.GameState, error) {
	registry := cases.NewCaseRegistry()

	casePath, err := registry.GetPath(caseID)
	if err != nil {
		return nil, fmt.Errorf("get case %q from registry: %w", caseID, err)
	}

	caseDef, err := cases.LoadCase(casePath)
	if err != nil {
		return nil, fmt.Errorf("load case %q: %w", caseID, err)
	}

	// clock := caseDef.StartTime{
	// 	CurrentTime: gameState.StartTime,
	// }
	// go BackgroundTime(game.Clock)

	// clock := game.Clock{
	// 	CurrentTime: game.StartTime,
	// },

	clock := game.NewClock(caseDef.StartTime.StartTime)

	gameState := &game.GameState{
		Case: caseDef.CaseInfo,

		Suspects:   caseDef.Suspects,
		Clues:      caseDef.Clues,
		Locations:  caseDef.Locations,
		Objects:    caseDef.Objects,
		Dialogues:  caseDef.Dialogues,
		Deductions: caseDef.Deductions,
		Facts:      caseDef.Facts,
		Motives:    caseDef.Motives,
		Events:     caseDef.Events,

		Clock:               clock,
		FoundClues:          make(map[int]bool),
		SearchedObjects:     make(map[int]bool),
		AskedDialogues:      make(map[game.DialogueID]bool),
		SolvedDeductions:    make(map[int]bool),
		FoundFacts:          make(map[int]bool),
		FoundMotives:        make(map[int]bool),
		ConfrontedDialogues: make(map[game.DialogueID]bool),

		OpenDialogues: game.InitOpenDialogues(caseDef.Dialogues),
		OpenLocations: game.InitOpenLocations(caseDef.Locations),
		OpenObjects:   game.InitOpenObjects(caseDef.Objects),
	}

	go clock.RunTime()

	return gameState, nil
}

func getStartedGame(w http.ResponseWriter, r *http.Request) (*Session, error) {
	session, err := getGame(w, r)
	if err != nil {
		return nil, err
	}

	session.Mu.RLock()
	defer session.Mu.RUnlock()

	if session.State == nil {
		return nil, ErrGameNotStarted
	}

	return session, nil
}

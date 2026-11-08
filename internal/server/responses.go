package server

import "detective/internal/game"

type StartResponse struct {
	Status    string `json:"status"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Intro     string `json:"intro"`
	KnownInfo string `json:"knownInfo"`
}

type InspectResponse struct {
	Object           game.Object `json:"object"`
	ClueFound        bool        `json:"clueFound"`
	ClueName         string      `json:"clueName"`
	ClueAbout        string      `json:"clueAbout"`
	DialogueUnlocked bool        `json:"dialogueUnlocked"`
}

type CaseResponse struct {
	Intro           string          `json:"intro"`
	KnownInfo       string          `json:"knownInfo"`
	FoundClues      []game.Clue     `json:"foundClues"`
	KnownFacts      []game.Dialogue `json:"knownFacts"`
	CluesFound      int             `json:"cluesFound"`
	CluesTotal      int             `json:"cluesTotal"`
	ObjectsSearched int             `json:"objectsSearched"`
	ObjectsTotal    int             `json:"objectsTotal"`
	DialoguesAsked  int             `json:"dialoguesAsked"`
	DialoguesTotal  int             `json:"dialoguesTotal"`
}

type AccuseResponse struct {
	TypeEnd string `json:"typeEnd"`
	TextEnd string `json:"textEnd"`
}

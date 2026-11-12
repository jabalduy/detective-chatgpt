package server

import (
 "detective/internal/game"
 "net/http"
)

// requireInvestigationAct prevents pre-discovery evidence disclosure and
// mutation. This is intentionally separate from the real-time deadline gate:
// old sessions with no narrative act retain investigation access.
func requireInvestigationAct(w http.ResponseWriter, state *game.GameState) bool {
 if state.CurrentNarrativeAct() == game.ActInvestigation {return true}
 http.Error(w,"investigation is not available during prologue",http.StatusConflict)
 return false
}

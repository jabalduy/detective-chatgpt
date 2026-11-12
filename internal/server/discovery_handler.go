package server

import (
 "detective/internal/game"
 "encoding/json"
 "errors"
 "net/http"
 "time"
)

// discoveryHandler is deliberately inert for legacy investigation sessions.
// Prologue activation and evidence gating will be added together in a later
// stage; this route cannot start a prologue or bypass its clock.
func discoveryHandler(w http.ResponseWriter,r *http.Request) {
 if r.Method!=http.MethodPost {http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return}
 session,err:=getGame(w,r)
 if err!=nil {http.Error(w,"session error",http.StatusInternalServerError);return}
 session.Mu.Lock()
 defer session.Mu.Unlock()
 if session.State==nil {http.Error(w,"game not started",http.StatusConflict);return}
 if session.State.CurrentNarrativeAct()!=game.ActPrologue {
  http.Error(w,"no active prologue",http.StatusConflict);return
 }
 changed,err:=session.State.BeginInvestigation(time.Now())
 if errors.Is(err,game.ErrDiscoveryTooEarly) || errors.Is(err,game.ErrDiscoveryInteractionActive) {
  http.Error(w,err.Error(),http.StatusConflict);return
 }
 if err!=nil {http.Error(w,err.Error(),http.StatusConflict);return}
 if !changed {http.Error(w,"discovery already completed",http.StatusConflict);return}
 // This handler is wrapped by saveAfterMutation, which persists the
 // successful transition with the same session's existing save mechanism.
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(map[string]string{"status":"investigation_started"})
}

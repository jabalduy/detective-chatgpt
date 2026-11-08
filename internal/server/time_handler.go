package server

import (
 "encoding/json"
 "net/http"
)

func timeHandler(w http.ResponseWriter, r *http.Request) {
 if r.Method != http.MethodGet {
  http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
  return
 }
 session, err := getStartedGame(w, r)
 if err != nil {
  http.Error(w, "game not started", http.StatusBadRequest)
  return
 }
 session.Mu.RLock()
 defer session.Mu.RUnlock()
 if session.State == nil || session.State.Clock == nil {
  http.Error(w, "game not started", http.StatusBadRequest)
  return
 }
 w.Header().Set("Content-Type", "application/json")
 _ = json.NewEncoder(w).Encode(struct {
  CurrentTime string `json:"current_time"`
 }{CurrentTime: session.State.Clock.GetTime().Format("15:04:05")})
}

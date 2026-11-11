package server

import (
 "net/http"
 "time"
)

// investigationDeadlineGuard covers legacy mutation routes as well as newer
// handlers. Reading evidence remains available during the final accusation.
// Ending an already-started conversation is explicitly exempt.
func investigationDeadlineGuard(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Method!=http.MethodPost {next.ServeHTTP(w,r);return}
  switch r.URL.Path {
  case "/scene/choose", "/conversation/start", "/examinations/start",
   "/evidence/present", "/inspect-object", "/ask-dialogue",
   "/dialogue/start", "/confront", "/deduce", "/event/choose", "/hint":
  default:
   next.ServeHTTP(w,r);return
  }
  session,err:=getGame(w,r)
  if err!=nil {http.Error(w,"session error",500);return}
  session.Mu.RLock()
  state:=session.State
  allowed:=state!=nil && state.CanInvestigateAt(time.Now())
  session.Mu.RUnlock()
  if state==nil {http.Error(w,"game not started",409);return}
  if !allowed {http.Error(w,"investigation time expired",409);return}
  next.ServeHTTP(w,r)
 })
}

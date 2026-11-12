package server

import (
 "encoding/json"
 "net/http"
 "detective/internal/game"
)

// Prologue travel deliberately never executes investigation scene actions.
func prologueTravelHandler(w http.ResponseWriter,r *http.Request) {
 if r.Method!=http.MethodPost {http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return}
 session,err:=getStartedGame(w,r)
 if err!=nil {http.Error(w,"game not started",http.StatusConflict);return}
 session.Mu.Lock()
 defer session.Mu.Unlock()
 s:=session.State
 if s.CurrentNarrativeAct()!=game.ActPrologue {http.Error(w,"prologue has ended",http.StatusConflict);return}
 if s.ActiveEvent!=nil || s.ActiveConversation!="" || s.ActiveDialogue!=nil {http.Error(w,"finish active interaction first",http.StatusConflict);return}
 destination:=r.URL.Query().Get("destination")
 allowed:=false
 for _,id:=range []string{"corridor","salon","staff","archive","vestibule"} {if destination==id {allowed=true;break}}
 if !allowed || game.SceneByID(s,destination)==nil {http.Error(w,"destination unavailable",http.StatusConflict);return}
 s.CurrentScene=destination
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(map[string]string{"current_scene":destination})
}

package server

import (
 "detective/internal/game"
 "encoding/json"
 "net/http"
)

func hintHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodGet && r.Method!=http.MethodPost {http.Error(w,"method not allowed",405);return}
 session,err:=getGame(w,r)
 if err!=nil {http.Error(w,"session error",500);return}
 session.Mu.Lock()
 defer session.Mu.Unlock()
 if session.State==nil {http.Error(w,"game not started",409);return}
 var view game.HintView
 if r.Method==http.MethodGet {view,err=game.AvailableHint(session.State)
 }else{view,err=game.TakeHint(session.State)}
 if err!=nil {http.Error(w,err.Error(),409);return}
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(view)
}

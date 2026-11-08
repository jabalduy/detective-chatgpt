package server

import (
 "detective/internal/game"
 "encoding/json"
 "net/http"
)

func trainHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodGet {http.Error(w,"method not allowed",405);return}
 session,err:=getGame(w,r)
 if err!=nil {http.Error(w,"session error",500);return}
 session.Mu.Lock()
 defer session.Mu.Unlock()
 if session.State==nil {http.Error(w,"game not started",409);return}
 game.TickTrain(session.State)
 game.TickEvents(session.State)
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(game.TrainView(session.State))
}

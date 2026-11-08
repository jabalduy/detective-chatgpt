package server

import (
 "encoding/json"
 "net/http"
)

func endingHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodGet {http.Error(w,"method not allowed",405);return}
 session,err:=getGame(w,r)
 if err!=nil {http.Error(w,"session error",500);return}
 session.Mu.RLock()
 defer session.Mu.RUnlock()
 if session.State==nil {http.Error(w,"game not started",409);return}
 if session.State.Ending==nil {w.WriteHeader(http.StatusNoContent);return}
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(session.State.Ending)
}

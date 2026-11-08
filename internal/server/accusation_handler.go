package server

import (
 "detective/internal/game"
 "encoding/json"
 "net/http"
)

type AccusationRequest = game.Accusation

func accuseHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodPost {http.Error(w,"method not allowed",405);return}
 session,err:=getGame(w,r)
 if err!=nil {http.Error(w,"session error",500);return}
 session.Mu.Lock()
 defer session.Mu.Unlock()
 if session.State==nil {http.Error(w,"game not started",409);return}
 var request game.Accusation
 if err:=json.NewDecoder(http.MaxBytesReader(w,r.Body,1<<20)).Decode(&request);err!=nil {http.Error(w,"invalid accusation",400);return}
 ending,err:=game.SubmitAccusation(session.State,request)
 if err!=nil {http.Error(w,err.Error(),409);return}
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(map[string]any{"type_end":ending.Type,"text_end":ending.Text,"title":ending.Title})
}

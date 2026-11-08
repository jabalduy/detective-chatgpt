package server

import (
 "detective/internal/game"
 "encoding/json"
 "net/http"
)

func examinationsHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodGet {http.Error(w,"method not allowed",405);return}
 session,err:=getGame(w,r)
 if err!=nil {http.Error(w,"session error",http.StatusInternalServerError);return}
 session.Mu.Lock();defer session.Mu.Unlock()
 if session.State==nil {http.Error(w,"game not started",409);return}
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(game.ExaminationViews(session.State))
}
func examinationStartHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodPost {http.Error(w,"method not allowed",405);return}
 session,err:=getGame(w,r)
 if err!=nil {http.Error(w,"session error",http.StatusInternalServerError);return}
 session.Mu.Lock();defer session.Mu.Unlock()
 if session.State==nil {http.Error(w,"game not started",409);return}
 if session.State.ActiveEvent!=nil {http.Error(w,"event active",409);return}
 if err:=game.StartExamination(session.State,r.URL.Query().Get("id"));err!=nil {http.Error(w,err.Error(),409);return}
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(game.ExaminationViews(session.State))
}
func evidenceHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodGet {http.Error(w,"method not allowed",405);return}
 session,err:=getGame(w,r)
 if err!=nil {http.Error(w,"session error",http.StatusInternalServerError);return}
 session.Mu.RLock();defer session.Mu.RUnlock()
 if session.State==nil {http.Error(w,"game not started",409);return}
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(game.AvailableEvidence(session.State))
}
func evidencePresentHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodPost {http.Error(w,"method not allowed",405);return}
 session,err:=getGame(w,r)
 if err!=nil {http.Error(w,"session error",http.StatusInternalServerError);return}
 session.Mu.Lock();defer session.Mu.Unlock()
 if session.State==nil {http.Error(w,"game not started",409);return}
 response,err:=game.PresentEvidence(session.State,r.URL.Query().Get("id"))
 if err!=nil {http.Error(w,err.Error(),409);return}
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(map[string]string{"response":response})
}

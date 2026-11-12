package server

import (
 "detective/internal/game"
 "encoding/json"
 "net/http"
 "time"
)

func conversationStartHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodPost {http.Error(w,"method not allowed",405);return}
 session,err:=getGame(w,r)
 if err!=nil {http.Error(w,"session error",http.StatusInternalServerError);return}
 session.Mu.Lock();defer session.Mu.Unlock()
 if session.State==nil {http.Error(w,"game not started",409);return}
 if session.State.CurrentNarrativeAct()!=game.ActPrologue && !session.State.CanInvestigateAt(time.Now()) {http.Error(w,"investigation time expired",409);return}
 if err:=game.StartConversation(session.State,r.URL.Query().Get("id"));err!=nil {http.Error(w,err.Error(),409);return}
 writeConversation(w,session.State)
}
func conversationHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodGet {http.Error(w,"method not allowed",405);return}
 session,err:=getGame(w,r)
 if err!=nil {http.Error(w,"session error",http.StatusInternalServerError);return}
 session.Mu.RLock();defer session.Mu.RUnlock()
 if session.State==nil {http.Error(w,"game not started",409);return}
 writeConversation(w,session.State)
}
func conversationChooseHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodPost {http.Error(w,"method not allowed",405);return}
 session,err:=getGame(w,r)
 if err!=nil {http.Error(w,"session error",http.StatusInternalServerError);return}
 session.Mu.Lock();defer session.Mu.Unlock()
 if session.State==nil {http.Error(w,"game not started",409);return}
 ended,err:=game.ChooseConversationOption(session.State,r.URL.Query().Get("choiceID"))
 if err!=nil {http.Error(w,err.Error(),409);return}
 if ended {session.State.FinishConversationGrace(time.Now());w.Header().Set("Content-Type","application/json");_=json.NewEncoder(w).Encode(map[string]bool{"ended":true});return}
 writeConversation(w,session.State)
}
func writeConversation(w http.ResponseWriter,s *game.GameState){
 view,err:=game.CurrentConversationView(s)
 if err!=nil {http.Error(w,err.Error(),409);return}
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(view)
}
func charactersHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodGet {http.Error(w,"method not allowed",405);return}
 session,err:=getGame(w,r)
 if err!=nil {http.Error(w,"session error",http.StatusInternalServerError);return}
 session.Mu.RLock();defer session.Mu.RUnlock()
 if session.State==nil {http.Error(w,"game not started",409);return}
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(session.State.Characters)
}

package server

import (
 "detective/internal/game"
 "encoding/json"
 "errors"
 "net/http"
 "os"
)

func saveHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodPost {http.Error(w,"method not allowed",405);return}
 session,err:=getGame(w,r);if err!=nil{http.Error(w,"session error",500);return}
 session.Mu.Lock();defer session.Mu.Unlock()
 if session.State==nil{http.Error(w,"game not started",409);return}
 id:=getSessionID(w,r)
 if err:=writeSave(savePath(id),snapshot(session.State));err!=nil{http.Error(w,"save failed",500);return}
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(map[string]string{"status":"saved"})
}

func loadHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodPost {http.Error(w,"method not allowed",405);return}
 session,err:=getGame(w,r);if err!=nil{http.Error(w,"session error",500);return}
 id:=getSessionID(w,r)
 session.Mu.Lock();defer session.Mu.Unlock()
 saved,err:=readSave(savePath(id))
 if r.URL.Query().Get("recovery")=="1" {
  checkpoint,checkpointErr:=readSave(recoveryPath(id))
  if checkpointErr==nil {saved,err=checkpoint,nil} else if !errors.Is(checkpointErr,os.ErrNotExist) {http.Error(w,"invalid recovery checkpoint",422);return}
 }
 if errors.Is(err,os.ErrNotExist){http.Error(w,"save not found",404);return}
 if err!=nil{http.Error(w,"invalid save",422);return}
 state,err:=newGame(saved.CaseID)
 if err!=nil{http.Error(w,"case unavailable",422);return}
 if session.State!=nil && session.State.Clock!=nil{session.State.Clock.StopTime()}
 applySave(state,saved)
 session.State=state
 game.TickTrain(state)
 game.TickEvents(state)
 // Loading may catch up scheduled train/event transitions. Persist those
 // changes immediately so another restart cannot replay them.
 if err:=writeSave(savePath(id),snapshot(state));err!=nil{
  http.Error(w,"save after load failed",http.StatusInternalServerError)
  return
 }
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(map[string]string{"status":"loaded","case_id":saved.CaseID})
}

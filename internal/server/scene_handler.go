package server

import (
 "encoding/json"
 "net/http"
 "time"
 "detective/internal/game"
)

func sceneHandler(w http.ResponseWriter,r *http.Request) {
 if r.Method != http.MethodGet { http.Error(w,"method not allowed",405); return }
 session,err:=getGame(w,r)
 if err!=nil { http.Error(w,"session error",500);return }
 session.Mu.RLock()
 defer session.Mu.RUnlock()
 if session.State==nil { http.Error(w,"game not started",409);return }
 if len(session.State.Scenes)==0 { http.Error(w,"no scenes defined",404);return }
 view,err:=game.CurrentSceneView(session.State)
 if err!=nil { http.Error(w,err.Error(),500);return }
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(view)
}

func sceneChooseHandler(w http.ResponseWriter,r *http.Request) {
 if r.Method!=http.MethodPost { http.Error(w,"method not allowed",405);return }
 session,err:=getGame(w,r)
 if err!=nil { http.Error(w,"session error",500);return }
 session.Mu.Lock()
 defer session.Mu.Unlock()
 if session.State==nil { http.Error(w,"game not started",409);return }
 if !session.State.CanInvestigateAt(time.Now()) {http.Error(w,"investigation time expired",409);return}
 if session.State.ActiveEvent!=nil || session.State.ActiveDialogue!=nil || session.State.ActiveConversation!="" {
  http.Error(w,"finish active event or dialogue first",409);return
 }
 id:=r.URL.Query().Get("choiceID")
 if id=="" { http.Error(w,"choiceID required",400);return }
 if err:=game.ChooseSceneOption(session.State,id);err!=nil {
  http.Error(w,err.Error(),409);return
 }
 view,err:=game.CurrentSceneView(session.State)
 if err!=nil { http.Error(w,err.Error(),500);return }
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(view)
}

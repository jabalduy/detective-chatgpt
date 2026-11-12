package server

import (
 "net/http"
 "testing"

 "detective/internal/game"
)

func TestPrologueDoesNotExposeInvestigationScenes(t *testing.T) {
 handler,cookie,session:=testSession(t)
 session.Mu.Lock()
 session.State=&game.GameState{
  NarrativeAct:game.ActPrologue,
  Scenes:[]game.Scene{{ID:"corridor",Title:"After the murder",Description:"A body was found"}},
  CurrentScene:"corridor",
 }
 session.Mu.Unlock()

 response:=requestWithCookie(handler,cookie,http.MethodGet,"/scene",nil)
 if response.Code!=http.StatusConflict {
  t.Fatalf("GET /scene status=%d want 409: %s",response.Code,response.Body.String())
 }
 session.Mu.RLock()
 view,err:=buildGameStateView(session.State)
 session.Mu.RUnlock()
 if err!=nil {t.Fatal(err)}
 if view.Scene!=nil {t.Fatal("prologue state exposed investigation scene")}
}

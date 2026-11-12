package server

import (
 "net/http"
 "testing"

 "detective/internal/game"
)

func TestLegacyDialogueRoutesRejectPrologue(t *testing.T) {
 handler,cookie,session:=testSession(t)
 session.Mu.Lock()
 session.State=&game.GameState{NarrativeAct:game.ActPrologue}
 session.Mu.Unlock()
 for _,tc:=range []struct{method,path string}{
  {http.MethodGet,"/dialogues?susID=1"},
  {http.MethodPost,"/askDialogue?susID=1&dialID=1"},
 } {
  r:=requestWithCookie(handler,cookie,tc.method,tc.path,nil)
  if r.Code!=http.StatusConflict {
   t.Fatalf("%s %s: status=%d want 409; body=%s",tc.method,tc.path,r.Code,r.Body.String())
  }
 }
}

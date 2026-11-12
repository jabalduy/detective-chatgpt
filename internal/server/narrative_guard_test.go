package server

import (
 "net/http"
 "net/http/httptest"
 "testing"

 "detective/internal/game"
)

func TestInvestigationGuardBlocksPrologueAndPreservesLegacy(t *testing.T) {
 for _,tc:=range []struct{name string;act game.NarrativeAct;allowed bool}{
  {"prologue",game.ActPrologue,false},
  {"investigation",game.ActInvestigation,true},
  {"legacy","",true},
  {"finished",game.ActFinished,false},
 } {
  t.Run(tc.name,func(t *testing.T){
   w:=httptest.NewRecorder()
   got:=requireInvestigationAct(w,&game.GameState{NarrativeAct:tc.act})
   if got!=tc.allowed {t.Fatalf("allowed=%v, want %v",got,tc.allowed)}
   if !tc.allowed && w.Code!=http.StatusConflict {t.Fatalf("status=%d",w.Code)}
  })
 }
}

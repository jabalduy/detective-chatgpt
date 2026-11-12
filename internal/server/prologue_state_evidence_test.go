package server

import (
 "encoding/json"
 "net/http"
 "testing"

 "detective/internal/game"
)

func TestPrologueStateDoesNotExposeInvestigationEvidenceIDs(t *testing.T) {
 handler,cookie,session:=testSession(t)
 session.Mu.Lock()
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"case_express_1928"},
  NarrativeAct:game.ActPrologue,
  FoundClues:map[int]bool{1:true},
  FoundFacts:map[int]bool{2:true},
  FoundMotives:map[int]bool{3:true},
 }
 session.Mu.Unlock()

 response:=requestWithCookie(handler,cookie,http.MethodGet,"/state",nil)
 if response.Code!=http.StatusOK {
  t.Fatalf("GET /state status=%d: %s",response.Code,response.Body.String())
 }
 var view gameStateView
 if err:=json.Unmarshal(response.Body.Bytes(),&view);err!=nil {t.Fatal(err)}
 if view.NarrativeAct!=game.ActPrologue {t.Fatalf("act=%s",view.NarrativeAct)}
 if len(view.FoundClueIDs)!=0 || len(view.FoundFactIDs)!=0 || len(view.FoundMotiveIDs)!=0 {
  t.Fatalf("prologue leaked evidence IDs: clues=%v facts=%v motives=%v",view.FoundClueIDs,view.FoundFactIDs,view.FoundMotiveIDs)
 }
 session.Mu.RLock()
 defer session.Mu.RUnlock()
 if !session.State.FoundClues[1] || !session.State.FoundFacts[2] || !session.State.FoundMotives[3] {
  t.Fatal("snapshot changed stored evidence")
 }
}

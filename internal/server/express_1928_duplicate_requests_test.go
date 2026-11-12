package server

import (
 "net/http"
 "path/filepath"
 "testing"

 "detective/internal/game"
)

// Operator double-clicks and stale retries must not replay a scene action
// after navigation, or award evidence when no interaction is active.
func TestExpress1928RepeatedSceneAndConversationRequests(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 t.Setenv("GAME_CASE_DIR",filepath.Join("..","..","cases"))
 handler,cookie,session:=testSession(t)
 check:=func(method,path string,want int){
  t.Helper()
  r:=requestWithCookie(handler,cookie,method,path,nil)
  if r.Code!=want {t.Fatalf("%s %s returned %d, want %d: %s",method,path,r.Code,want,r.Body.String())}
 }
 check(http.MethodPost,"/start?caseID=case_express_1928",200)
 defer session.State.Clock.StopTime()
 check(http.MethodPost,"/scene/choose?choiceID=victor",200)
 check(http.MethodPost,"/scene/choose?choiceID=victor",409)
 if session.State.CurrentScene!="compartment" {t.Fatal("stale navigation changed location")}
 check(http.MethodPost,"/scene/choose?choiceID=ribbon",200)
 check(http.MethodPost,"/scene/choose?choiceID=back",200)
 check(http.MethodPost,"/scene/choose?choiceID=salon",200)
 check(http.MethodPost,"/conversation/start?id=boris_intro",200)
 check(http.MethodPost,"/conversation/start?id=boris_intro",409)
 check(http.MethodPost,"/scene/choose?choiceID=back",409)
 check(http.MethodPost,"/conversation/choose?choiceID=leave",200)
 check(http.MethodPost,"/conversation/choose?choiceID=leave",409)
 if !session.State.FoundFacts[1] {t.Fatal("ribbon evidence lost")}
 if len(session.State.FoundMotives)!=0 || session.State.Ending!=nil {t.Fatal("retry mutated verdict or motives")}
}

// A finished verdict must remain immutable even when an operator sends
// an outdated action or repeats an accusation.
func TestExpress1928FinishedGameRejectsStaleMutations(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 t.Setenv("GAME_CASE_DIR",filepath.Join("..","..","cases"))
 handler,cookie,session:=testSession(t)
 check:=func(method,path string,want int){
  t.Helper()
  r:=requestWithCookie(handler,cookie,method,path,nil)
  if r.Code!=want {t.Fatalf("%s %s returned %d, want %d: %s",method,path,r.Code,want,r.Body.String())}
 }
 check(http.MethodPost,"/start?caseID=case_express_1928",200)
 defer session.State.Clock.StopTime()
 session.State.Ending=&game.GameEnding{Type:"win",Title:"Finished",Text:"Resolved"}
 check(http.MethodPost,"/scene/choose?choiceID=victor",409)
 check(http.MethodPost,"/conversation/start?id=boris_intro",409)
 if session.State.CurrentScene!="corridor" || session.State.Ending.Title!="Finished" {
  t.Fatal("stale request mutated completed game")
 }
}

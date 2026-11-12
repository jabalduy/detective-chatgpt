package server

import (
 "encoding/json"
 "net/http"
 "path/filepath"
 "testing"
 "time"

 "detective/internal/game"
)

// A rooftop decision is a real HTTP event, and reloading cannot reset its
// deadline or lose the transition from the first timed decision.
func TestExpress1928RooftopChaseHTTPAndRestart(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 t.Setenv("GAME_CASE_DIR",filepath.Join("..","..","cases"))
 handler,cookie,session:=testSession(t)
 check:=func(method,path string,want int)[]byte{
  t.Helper()
  r:=requestWithCookie(handler,cookie,method,path,nil)
  if r.Code!=want {t.Fatalf("%s %s: got %d want %d: %s",method,path,r.Code,want,r.Body.String())}
  return r.Body.Bytes()
 }
 check(http.MethodPost,"/start?caseID=case_express_1928",200)
 defer func(){if session.State!=nil && session.State.Clock!=nil {session.State.Clock.StopTime()}}()
 session.Mu.Lock()
 if !game.TriggerEvent(session.State,103) {session.Mu.Unlock();t.Fatal("could not trigger chase")}
 session.Mu.Unlock()
 var first EventResponse
 if err:=json.Unmarshal(check(http.MethodGet,"/event",200),&first);err!=nil {t.Fatal(err)}
 if first.ID!=103 {t.Fatalf("first event=%d",first.ID)}
 check(http.MethodPost,"/event/choose?optionID=1",200)
 var roof EventResponse
 if err:=json.Unmarshal(check(http.MethodGet,"/event",200),&roof);err!=nil {t.Fatal(err)}
 if roof.ID!=105 {t.Fatalf("expected rooftop follow-up, got %d",roof.ID)}
 session.Mu.RLock()
 deadline:=session.State.EventDeadline
 session.Mu.RUnlock()
 check(http.MethodPost,"/save",200)
 check(http.MethodPost,"/load",200)
 session.Mu.RLock()
 if session.State.ActiveEvent==nil || *session.State.ActiveEvent!=105 || !session.State.EventDeadline.Equal(deadline) {
  session.Mu.RUnlock();t.Fatal("reload lost rooftop event or reset countdown")
 }
 session.Mu.RUnlock()
 check(http.MethodPost,"/event/choose?optionID=1",200)
 session.Mu.RLock()
 finished:=session.State.ActiveEvent==nil
 evidence:=session.State.FoundFacts[18]
 session.Mu.RUnlock()
 if !finished || !evidence {t.Fatal("rooftop outcome missing after reload")}
 check(http.MethodPost,"/event/choose?optionID=1",409)
}

func TestExpress1928RooftopTimeoutViaHTTP(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 t.Setenv("GAME_CASE_DIR",filepath.Join("..","..","cases"))
 handler,cookie,session:=testSession(t)
 r:=requestWithCookie(handler,cookie,http.MethodPost,"/start?caseID=case_express_1928",nil)
 if r.Code!=200 {t.Fatalf("start: %d",r.Code)}
 defer session.State.Clock.StopTime()
 session.Mu.Lock()
 if !game.TriggerEvent(session.State,105) {session.Mu.Unlock();t.Fatal("cannot trigger rooftop event")}
 session.State.EventDeadline=time.Now().Add(-time.Second)
 session.Mu.Unlock()
 r=requestWithCookie(handler,cookie,http.MethodPost,"/event/expire",nil)
 if r.Code!=204 {t.Fatalf("expire: %d %s",r.Code,r.Body.String())}
 session.Mu.RLock()
 active:=session.State.ActiveEvent!=nil
 found:=session.State.FoundFacts[18]
 session.Mu.RUnlock()
 if active || found {t.Fatal("timeout failed to retreat without granting evidence")}
}

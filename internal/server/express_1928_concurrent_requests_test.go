package server

import (
 "net/http"
 "path/filepath"
 "sync"
 "testing"
)

// Concurrent requests share one operator session. Only one stale navigation
// can succeed, and the other requests must leave the same coherent state.
func TestExpress1928ConcurrentNavigationIsSerialized(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 t.Setenv("GAME_CASE_DIR",filepath.Join("..","..","cases"))
 handler,cookie,session:=testSession(t)
 start:=requestWithCookie(handler,cookie,http.MethodPost,"/start?caseID=case_express_1928",nil)
 if start.Code!=200 {t.Fatalf("start: %d %s",start.Code,start.Body.String())}
 defer session.State.Clock.StopTime()
 const workers=12
 statuses:=make(chan int,workers)
 var wg sync.WaitGroup
 for i:=0;i<workers;i++ {
  wg.Add(1)
  go func(){
   defer wg.Done()
   response:=requestWithCookie(handler,cookie,http.MethodPost,"/scene/choose?choiceID=victor",nil)
   statuses<-response.Code
  }()
 }
 wg.Wait()
 close(statuses)
 success,conflict:=0,0
 for status:=range statuses {
  switch status {
  case 200:success++
  case 409:conflict++
  default:t.Errorf("unexpected HTTP status %d",status)
  }
 }
 if success!=1 || conflict!=workers-1 {t.Fatalf("success=%d conflict=%d",success,conflict)}
 session.Mu.RLock()
 scene:=session.State.CurrentScene
 facts:=len(session.State.FoundFacts)
 session.Mu.RUnlock()
 if scene!="compartment" || facts!=0 {t.Fatalf("state corrupted: scene=%s facts=%d",scene,facts)}
 saved,err:=readSave(savePath(cookie.Value))
 if err!=nil {t.Fatal(err)}
 if saved.CurrentScene!="compartment" {t.Fatalf("autosave scene=%q",saved.CurrentScene)}
}

func TestExpress1928ConcurrentConversationStartIsSerialized(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 t.Setenv("GAME_CASE_DIR",filepath.Join("..","..","cases"))
 handler,cookie,session:=testSession(t)
 for _,path:=range []string{
  "/start?caseID=case_express_1928",
  "/scene/choose?choiceID=salon",
 } {
  response:=requestWithCookie(handler,cookie,http.MethodPost,path,nil)
  if response.Code!=200 {t.Fatalf("%s: %d %s",path,response.Code,response.Body.String())}
 }
 defer session.State.Clock.StopTime()
 const workers=10
 statuses:=make(chan int,workers)
 var wg sync.WaitGroup
 for i:=0;i<workers;i++ {
  wg.Add(1)
  go func(){
   defer wg.Done()
   response:=requestWithCookie(handler,cookie,http.MethodPost,"/conversation/start?id=boris_intro",nil)
   statuses<-response.Code
  }()
 }
 wg.Wait()
 close(statuses)
 success,conflict:=0,0
 for status:=range statuses {
  switch status {
  case 200:success++
  case 409:conflict++
  default:t.Errorf("unexpected HTTP status %d",status)
  }
 }
 if success!=1 || conflict!=workers-1 {t.Fatalf("success=%d conflict=%d",success,conflict)}
 session.Mu.RLock()
 active:=session.State.ActiveConversation
 paused:=session.State.Clock.IsPaused()
 session.Mu.RUnlock()
 if active!="boris_intro" || !paused {t.Fatalf("conversation=%q paused=%t",active,paused)}
 saved,err:=readSave(savePath(cookie.Value))
 if err!=nil {t.Fatal(err)}
 if saved.ActiveConversation!="boris_intro" || !saved.ClockPaused {
  t.Fatal("autosave lost active conversation or clock pause")
 }
}

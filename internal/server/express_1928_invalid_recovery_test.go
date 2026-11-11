package server

import (
 "net/http"
 "os"
 "path/filepath"
 "testing"
 "time"

 "detective/internal/game"
)

// Failed restores must never replace a running investigation or its clock.
func TestCorruptLoadPreservesActiveInvestigation(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(1928,10,1,21,0,0,0,time.UTC)
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"case_express_1928"},
  Clock:game.NewClock(start),
  CurrentScene:"staff",
  FoundFacts:map[int]bool{1:true,3:true},
 }
 defer session.State.Clock.StopTime()
 original:=session.State
 for _,tc:=range []struct{name,body string;want int}{
  {"malformed json","{not-json",422},
  {"unsupported version",`{"version":999,"case_id":"case_express_1928","clock_time":"1928-10-01T21:00:00Z"}`,422},
  {"missing case",`{"version":1,"clock_time":"1928-10-01T21:00:00Z"}`,422},
  {"missing clock",`{"version":1,"case_id":"case_express_1928"}`,422},
  {"unknown case",`{"version":1,"case_id":"nonexistent_case","clock_time":"1928-10-01T21:00:00Z"}`,422},
 } {
  t.Run(tc.name,func(t *testing.T){
   if err:=os.MkdirAll(filepath.Dir(savePath(cookie.Value)),0700);err!=nil {t.Fatal(err)}
   if err:=os.WriteFile(savePath(cookie.Value),[]byte(tc.body),0600);err!=nil {t.Fatal(err)}
   result:=requestWithCookie(handler,cookie,http.MethodPost,"/load",nil)
   if result.Code!=tc.want {t.Fatalf("load status=%d want=%d: %s",result.Code,tc.want,result.Body.String())}
   if session.State!=original {t.Fatal("failed load replaced current state")}
   if session.State.CurrentScene!="staff" || !session.State.FoundFacts[1] || !session.State.FoundFacts[3] {
    t.Fatal("failed load discarded investigation progress")
   }
   if !session.State.Clock.GetTime().Equal(start) {t.Fatal("failed load reset route clock")}
  })
 }
}

func TestInvalidExpressActionsDoNotAwardEvidence(t *testing.T){
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 state,err:=newGame("case_express_1928")
 if err!=nil {t.Fatal(err)}
 session.State=state
 defer state.Clock.StopTime()
 originalScene:=state.CurrentScene
 for _,tc:=range []struct{method,path string;want int}{
  {http.MethodGet,"/scene/choose?choiceID=victor",405},
  {http.MethodPost,"/scene/choose?choiceID=unknown",409},
  {http.MethodPost,"/scene/choose?choiceID=box",409},
  {http.MethodPost,"/conversation/start?id=boris_intro",409},
  {http.MethodPost,"/conversation/choose?choiceID=clock_confront",409},
 } {
  r:=requestWithCookie(handler,cookie,tc.method,tc.path,nil)
  if r.Code!=tc.want {t.Fatalf("%s %s: status=%d want=%d: %s",tc.method,tc.path,r.Code,tc.want,r.Body.String())}
  if state.CurrentScene!=originalScene || len(state.FoundFacts)!=0 || len(state.FoundMotives)!=0 {
   t.Fatalf("rejected request changed evidence or scene: %s",tc.path)
  }
 }
}

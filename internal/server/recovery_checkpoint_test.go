package server

import (
 "net/http"
 "path/filepath"
 "detective/internal/cases"
 "testing"
 "time"

 "detective/internal/game"
)

func TestRecoveryRestoresPreviousActionAndSurvivesStatePolling(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 def,err:=cases.LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 state:=&game.GameState{Case:def.CaseInfo,Clock:game.NewClock(time.Date(1928,10,1,21,0,0,0,time.UTC)),Scenes:def.Scenes,Conversations:def.Conversations,Characters:def.Characters,CurrentScene:"corridor",FoundFacts:map[int]bool{},FoundMotives:map[int]bool{},SolvedDeductions:map[int]bool{}}
 defer state.Clock.StopTime()
 session.Mu.Lock()
 session.State=state
 session.State.NarrativeAct=game.ActInvestigation
 session.State.InvestigationDeadline=time.Now().Add(time.Hour)
 session.State.AccusationDeadline=time.Now().Add(65*time.Minute)
 session.Mu.Unlock()
 check:=func(method,path string,want int) {
  t.Helper()
  res:=requestWithCookie(handler,cookie,method,path,nil)
  if res.Code!=want {t.Fatalf("%s %s: %d %s",method,path,res.Code,res.Body.String())}
 }
 check(http.MethodPost,"/save",200)
 check(http.MethodPost,"/scene/choose?choiceID=victor",200)
 check(http.MethodGet,"/state",200)
 check(http.MethodGet,"/state",200)
 saved,err:=readSave(recoveryPath(cookie.Value))
 if err!=nil {t.Fatal(err)}
 if saved.CurrentScene!="corridor" {t.Fatalf("checkpoint overwritten by polling: %q",saved.CurrentScene)}
 t.Chdir(filepath.Join("..",".."))
 check(http.MethodPost,"/load?recovery=1",200)
 session.Mu.RLock()
 scene:=session.State.CurrentScene
 session.Mu.RUnlock()
 if scene!="corridor" {t.Fatalf("restored scene=%q, want corridor",scene)}
}

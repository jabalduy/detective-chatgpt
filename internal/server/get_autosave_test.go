package server

import (
 "detective/internal/game"
 "net/http"
 "testing"
 "time"
)

func TestGETProgressIsAutosaved(t *testing.T){
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"test"},Clock:game.NewClock(start),RouteStart:start,
  Stations:[]game.Station{{ID:"stop",Arrival:"23:55",ArriveEffects:[]game.Effect{{Type:"discover_fact",ID:17}}}},
  FoundFacts:map[int]bool{},
 }
 defer session.State.Clock.StopTime()
 session.State.Clock.Add(6*time.Minute)
 r:=requestWithCookie(handler,cookie,http.MethodGet,"/train",nil)
 if r.Code!=http.StatusOK{t.Fatalf("GET /train status=%d: %s",r.Code,r.Body.String())}
 saved,err:=readSave(savePath(cookie.Value))
 if err!=nil{t.Fatal(err)}
 if !saved.FoundFacts[17] || !saved.ProcessedStations["stop:arrival"]{
  t.Fatalf("train progress not saved: %+v",saved)
 }
}

func TestGETExaminationCompletionIsAutosaved(t *testing.T){
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"test"},Clock:game.NewClock(start),
  Examinations:[]game.Examination{{ID:"lab",DurationSeconds:5,Effects:[]game.Effect{{Type:"discover_fact",ID:19}}}},
  ExaminationStarted:map[string]time.Time{"lab":start},
  FoundFacts:map[int]bool{},
 }
 defer session.State.Clock.StopTime()
 session.State.Clock.Add(6*time.Second)
 r:=requestWithCookie(handler,cookie,http.MethodGet,"/examinations",nil)
 if r.Code!=http.StatusOK{t.Fatalf("GET /examinations status=%d: %s",r.Code,r.Body.String())}
 saved,err:=readSave(savePath(cookie.Value))
 if err!=nil{t.Fatal(err)}
 if !saved.ExaminationCompleted["lab"] || !saved.FoundFacts[19]{
  t.Fatalf("examination completion not saved: %+v",saved)
 }
}

func TestAutosaveOnlyExpectedMethods(t *testing.T){
 for _,tc:=range []struct{method,path string;want bool}{
  {"GET","/train",true},{"GET","/event",true},{"GET","/examinations",true},
  {"GET","/state",true},{"GET","/scene",false},{"POST","/scene/choose",true},
  {"POST","/save",false},{"POST","/load",false},{"DELETE","/train",false},
 }{
  req,_:=http.NewRequest(tc.method,tc.path,nil)
  if got:=shouldAutosave(req);got!=tc.want{t.Errorf("%s %s: got %v want %v",tc.method,tc.path,got,tc.want)}
 }
}

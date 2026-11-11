package server

import (
 "net/http"
 "testing"
 "time"
 "detective/internal/game"
)

func TestInvestigationCommandsBlockedAfterDeadline(t *testing.T){
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,23,0,0,0,time.UTC)
 session.State=&game.GameState{
  Clock:game.NewClock(start),CurrentScene:"room",
  InvestigationDeadline:time.Now().Add(-time.Minute),
  AccusationDeadline:time.Now().Add(4*time.Minute),
  Scenes:[]game.Scene{{ID:"room",Choices:[]game.Choice{{ID:"look"}}}},
  Examinations:[]game.Examination{{ID:"lab"}},
  FoundClues:map[int]bool{},FoundFacts:map[int]bool{},
 }
 defer session.State.Clock.StopTime()
 for _,path:=range []string{
  "/scene/choose?choiceID=look",
  "/examinations/start?id=lab",
  "/evidence/present?id=proof",
  "/conversation/start?id=witness",
  "/inspect-object?objID=1",
  "/ask-dialogue?susID=1&dialID=1",
  "/dialogue/start?susID=1&dialID=1",
  "/confront?susID=1&dialID=1",
  "/deduce?id=1",
  "/event/choose?optionID=1",
  "/hint",
 }{
  response:=requestWithCookie(handler,cookie,http.MethodPost,path,nil)
  if response.Code!=http.StatusConflict {t.Errorf("%s returned %d, want 409",path,response.Code)}
 }
}

func TestExpiredInvestigationDoesNotConsumeHint(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,23,0,0,0,time.UTC)
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"test"},
  Clock:game.NewClock(start),
  InvestigationDeadline:time.Now().Add(-time.Minute),
  AccusationDeadline:time.Now().Add(4*time.Minute),
  Hints:[]game.Hint{},
 }
 defer session.State.Clock.StopTime()
 response:=requestWithCookie(handler,cookie,http.MethodPost,"/hint",nil)
 if response.Code!=http.StatusConflict {t.Fatalf("late hint returned %d",response.Code)}
 if session.State.HintCount!=0 || len(session.State.UsedHints)!=0 {
  t.Fatal("late hint changed game state")
 }
}

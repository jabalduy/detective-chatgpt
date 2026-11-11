package server

import (
 "detective/internal/game"
 "net/http"
 "testing"
 "time"
)

func TestQuestionBudgetAcrossDialogueAPIs(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 now:=time.Date(2026,1,1,12,0,0,0,time.UTC)
 first:=game.DialogueID{SusID:7,DialID:1}
 second:=game.DialogueID{SusID:7,DialID:2}
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"test"},Clock:game.NewClock(now),
  Suspects:[]game.Suspect{{SusID:7,MaxQuestions:1}},
  Dialogues:[]game.Dialogue{{SusID:7,DialID:1},{SusID:7,DialID:2}},
  OpenDialogues:map[game.DialogueID]bool{first:true,second:true},
  AskedDialogues:make(map[game.DialogueID]bool),
 }
 defer session.State.Clock.StopTime()
 req:=func(method,path string,want int) {
  t.Helper()
  r:=requestWithCookie(handler,cookie,method,path,nil)
  if r.Code!=want {t.Fatalf("%s %s: got %d want %d: %s",method,path,r.Code,want,r.Body.String())}
 }
 req(http.MethodPost,"/dialogue/start?susID=7&dialID=1",200)
 if session.State.AskedQuestionsBySuspect[7]!=0 {t.Fatal("starting a dialogue must not spend budget")}
 req(http.MethodPost,"/dialogue/finish?susID=7&dialID=1",200)
 if session.State.AskedQuestionsBySuspect[7]!=1 {t.Fatal("finishing must spend one question")}
 req(http.MethodPost,"/dialogue/start?susID=7&dialID=2",409)
 req(http.MethodPost,"/ask-dialogue?susID=7&dialID=2",409)
 req(http.MethodPost,"/ask-dialogue?susID=7&dialID=1",200)
 if session.State.AskedQuestionsBySuspect[7]!=1 {t.Fatal("rereading a question must not spend budget")}
 saved:=snapshot(session.State)
 if saved.AskedQuestionsBySuspect[7]!=1 {t.Fatal("question budget not saved")}
 restored:=&game.GameState{Clock:game.NewClock(now)}
 defer restored.Clock.StopTime()
 applySave(restored,saved)
 if restored.AskedQuestionsBySuspect[7]!=1 {t.Fatal("question budget not restored")}
}

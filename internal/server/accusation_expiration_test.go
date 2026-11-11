package server

import (
 "net/http"
 "testing"
 "time"
 "detective/internal/game"
)

func TestAccusationExpiresAndPersists(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,23,0,0,0,time.UTC)
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"timeout-case",Endings:game.Endings{Unsolved:"Case remains open"}},
  Clock:game.NewClock(start),
  InvestigationDeadline:time.Now().Add(-10*time.Minute),
  AccusationDeadline:time.Now().Add(-5*time.Minute),
 }
 defer session.State.Clock.StopTime()
 response:=requestWithCookie(handler,cookie,http.MethodGet,"/ending",nil)
 if response.Code!=200 {t.Fatalf("ending returned %d: %s",response.Code,response.Body.String())}
 if session.State.Ending==nil || session.State.Ending.Type!="unsolved" || !session.State.Clock.IsPaused(){
  t.Fatal("expired game was not finalized")
 }
 saved,err:=readSave(savePath(cookie.Value))
 if err!=nil || saved.Ending==nil || saved.Ending.Type!="unsolved"{
  t.Fatalf("timeout not persisted: %+v %v",saved.Ending,err)
 }
 if response:=requestWithCookie(handler,cookie,http.MethodPost,"/accuse",[]byte("{}"));response.Code!=409{
  t.Fatalf("late accusation returned %d",response.Code)
 }
}

func TestAccusationWaitsForConversationToFinish(t *testing.T){
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,23,0,0,0,time.UTC)
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"grace-case"},
  Clock:game.NewClock(start),
  InvestigationDeadline:time.Now().Add(-10*time.Minute),
  AccusationDeadline:time.Now().Add(-5*time.Minute),
  ActiveConversation:"witness",ConversationNode:"first",
  Conversations:[]game.Conversation{{ID:"witness",StartNode:"first",Nodes:[]game.ConversationNode{{ID:"first",Choices:[]game.ConversationChoice{{ID:"finish",End:true}}}}}},
 }
 defer session.State.Clock.StopTime()
 response:=requestWithCookie(handler,cookie,http.MethodGet,"/ending",nil)
 if response.Code!=204 || session.State.Ending!=nil {t.Fatal("conversation interrupted by timeout")}
 response=requestWithCookie(handler,cookie,http.MethodPost,"/conversation/choose?choiceID=finish",nil)
 if response.Code!=200 {t.Fatalf("could not finish conversation: %d %s",response.Code,response.Body.String())}
 if session.State.Ending!=nil {t.Fatal("conversation grace must not immediately expire the game")}
 if session.State.AccusationDeadline.Before(time.Now().Add(4*time.Minute)) {
  t.Fatal("conversation ending after deadline must grant five minutes")
 }
 if session.State.PhaseAt(time.Now())!=game.PhaseAccusation {t.Fatal("accusation phase not started")}
 saved,err:=readSave(savePath(cookie.Value))
 if err!=nil || !saved.AccusationDeadline.Equal(session.State.AccusationDeadline) {
  t.Fatalf("extended deadline not persisted: %v",err)
 }
}

package server

import (
 "detective/internal/game"
 "encoding/json"
 "net/http"
 "testing"
 "time"
)

func TestStatePollingResolvesExpiredEventAndAutosaves(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,23,0,0,0,time.UTC)
 id:=5
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"test"},Clock:game.NewClock(start),
  ActiveEvent:&id,EventDeadline:time.Now().Add(-time.Second),
  FoundFacts:map[int]bool{},
  Events:[]game.Event{{ID:5,TimeLimit:10,DefaultOptionID:20,Options:[]game.EventOption{
   {ID:10,Action:game.Action{Effects:[]game.Effect{{Type:"discover_fact",ID:10}}}},
   {ID:20,Action:game.Action{Effects:[]game.Effect{{Type:"discover_fact",ID:20}}}},
  }}},
 }
 defer session.State.Clock.StopTime()
 response:=requestWithCookie(handler,cookie,http.MethodGet,"/state",nil)
 if response.Code!=200 {t.Fatalf("state: %d %s",response.Code,response.Body.String())}
 var view gameStateView
 if err:=json.Unmarshal(response.Body.Bytes(),&view);err!=nil{t.Fatal(err)}
 if view.ActiveEvent!=nil || session.State.ActiveEvent!=nil {t.Fatal("expired event still active")}
 if !session.State.FoundFacts[20] || session.State.FoundFacts[10] {t.Fatal("wrong default outcome")}
 saved,err:=readSave(savePath(cookie.Value))
 if err!=nil || !saved.FoundFacts[20] || saved.ActiveEvent!=nil {
  t.Fatalf("default consequence not persisted: %+v %v",saved.ActiveEvent,err)
 }
}

func TestStatePollingTriggersScheduledEvent(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,23,0,0,0,time.UTC)
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"test"},Clock:game.NewClock(start),
  RouteStart:start,
  Events:[]game.Event{{ID:7,At:"23:00",TimeLimit:10,Options:[]game.EventOption{{ID:1,Text:"Respond"}}}},
 }
 defer session.State.Clock.StopTime()
 response:=requestWithCookie(handler,cookie,http.MethodGet,"/state",nil)
 if response.Code!=200 {t.Fatalf("state: %d %s",response.Code,response.Body.String())}
 if session.State.ActiveEvent==nil || *session.State.ActiveEvent!=7 {t.Fatal("scheduled event did not trigger")}
 saved,err:=readSave(savePath(cookie.Value))
 if err!=nil || saved.ActiveEvent==nil || *saved.ActiveEvent!=7 {
  t.Fatalf("triggered event not persisted: %v",err)
 }
}

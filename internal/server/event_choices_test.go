package server

import (
 "detective/internal/game"
 "encoding/json"
 "net/http"
 "testing"
 "time"
)

func TestEventAPIHidesLockedChoices(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,12,0,0,0,time.UTC)
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"test"},
  Clock:game.NewClock(start),
  FoundClues:make(map[int]bool),
  FoundFacts:make(map[int]bool),
  Events:[]game.Event{{ID:1,TimeLimit:60,Options:[]game.EventOption{
   {ID:10,Text:"Show evidence",Action:game.Action{Requirements:[]game.Requirement{{Type:"clue_found",ID:5}},Effects:[]game.Effect{{Type:"discover_fact",ID:10}}}},
   {ID:20,Text:"Wait",Action:game.Action{Effects:[]game.Effect{{Type:"discover_fact",ID:20}}}},
  }}},
 }
 defer session.State.Clock.StopTime()
 if !game.TriggerEvent(session.State,1){t.Fatal("trigger failed")}
 response:=requestWithCookie(handler,cookie,http.MethodGet,"/event",nil)
 if response.Code!=200 {t.Fatalf("event status %d: %s",response.Code,response.Body.String())}
 var view EventResponse
 if err:=json.Unmarshal(response.Body.Bytes(),&view);err!=nil{t.Fatal(err)}
 if len(view.Options)!=1 || view.Options[0].ID!=20 {t.Fatalf("locked option leaked: %+v",view.Options)}
 response=requestWithCookie(handler,cookie,http.MethodPost,"/event/choose?optionID=10",nil)
 if response.Code!=409 || session.State.FoundFacts[10] {t.Fatal("locked event choice accepted")}
 session.State.FoundClues[5]=true
 response=requestWithCookie(handler,cookie,http.MethodGet,"/event",nil)
 if err:=json.Unmarshal(response.Body.Bytes(),&view);err!=nil{t.Fatal(err)}
 if len(view.Options)!=2 {t.Fatalf("evidence did not unlock option: %+v",view.Options)}
 response=requestWithCookie(handler,cookie,http.MethodPost,"/event/choose?optionID=10",nil)
 if response.Code!=200 || !session.State.FoundFacts[10] || session.State.ActiveEvent!=nil {
  t.Fatalf("unlocked event choice failed: %d %s",response.Code,response.Body.String())
 }
}

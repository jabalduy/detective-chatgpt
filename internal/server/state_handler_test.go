package server

import (
 "detective/internal/game"
 "encoding/json"
 "net/http"
 "strings"
 "testing"
 "time"
)

func TestStateEndpointRequiresGame(t *testing.T){
 handler,cookie,_:=testSession(t)
 r:=requestWithCookie(handler,cookie,http.MethodGet,"/state",nil)
 if r.Code!=http.StatusConflict{t.Fatalf("status=%d body=%s",r.Code,r.Body.String())}
}

func TestStateEndpointIsReadOnlyAndHidesSpoilers(t *testing.T){
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 s:=&game.GameState{
  Case:game.CaseInfo{ID:"test",Solution:game.Solution{SuspectID:99}},
  Clock:game.NewClock(start),RouteStart:start,
  CurrentScene:"corridor",
  Scenes:[]game.Scene{{ID:"corridor",Choices:[]game.Choice{{ID:"search",Text:"Search",Action:game.Action{Effects:[]game.Effect{{Type:"discover_clue",ID:44}}}}}}},
  FoundClues:map[int]bool{3:true,1:true},FoundFacts:map[int]bool{2:true},FoundMotives:map[int]bool{},
  Characters:map[int]game.CharacterState{1:{Location:"corridor",Present:true}},
  Events:[]game.Event{{ID:9,Title:"Secret event",Options:[]game.EventOption{{ID:5,Text:"Choose",Action:game.Action{Effects:[]game.Effect{{Type:"discover_fact",ID:88}}}}}}},
 }
 session.State=s
 defer s.Clock.StopTime()
 r:=requestWithCookie(handler,cookie,http.MethodGet,"/state",nil)
 if r.Code!=http.StatusOK{t.Fatalf("status=%d body=%s",r.Code,r.Body.String())}
 body:=r.Body.String()
 for _,secret:=range []string{"suspect_id","discover_clue","discover_fact","Secret event"}{
  if strings.Contains(body,secret){t.Fatalf("spoiler %q leaked: %s",secret,body)}
 }
 var view gameStateView
 if err:=json.Unmarshal(r.Body.Bytes(),&view);err!=nil{t.Fatal(err)}
 if view.Scene==nil || view.Scene.ID!="corridor" || len(view.FoundClueIDs)!=2 || view.FoundClueIDs[0]!=1 || view.FoundClueIDs[1]!=3{
  t.Fatalf("unexpected state: %+v",view)
 }
 if s.FoundClues[44]{t.Fatal("GET /state mutated game state")}
}

func TestStateEndpointShowsActiveEventWithoutEffects(t *testing.T){
 handler,cookie,session:=testSession(t)
 id:=9
 session.State=&game.GameState{
  ActiveEvent:&id,EventDeadline:time.Now().Add(20*time.Second),
  Events:[]game.Event{{ID:9,Title:"Emergency",Options:[]game.EventOption{{ID:5,Text:"Help",Action:game.Action{Effects:[]game.Effect{{Type:"discover_fact",ID:77}}}}}}},
 }
 r:=requestWithCookie(handler,cookie,http.MethodGet,"/state",nil)
 if r.Code!=http.StatusOK{t.Fatalf("status=%d body=%s",r.Code,r.Body.String())}
 var view gameStateView
 if err:=json.Unmarshal(r.Body.Bytes(),&view);err!=nil{t.Fatal(err)}
 if view.ActiveEvent==nil || view.ActiveEvent.ID!=9 || len(view.ActiveEvent.Options)!=1 || view.ActiveEvent.RemainingSeconds<1{
  t.Fatalf("active event missing: %+v",view.ActiveEvent)
 }
 if strings.Contains(r.Body.String(),"discover_fact"){t.Fatal("event effect leaked")}
}

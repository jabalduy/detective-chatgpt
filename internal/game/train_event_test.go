package game

import (
 "testing"
 "time"
)

func TestTrainStationEffectsAfterTimeJump(t *testing.T){
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 s:=&GameState{Clock:NewClock(start),RouteStart:start,
  FoundFacts:make(map[int]bool),
  Stations:[]Station{{ID:"one",Name:"Night stop",Arrival:"00:10",Departure:"00:12",ArriveEffects:[]Effect{{Type:"discover_fact",ID:7}}}},
 }
 s.Clock.Add(21*time.Minute)
 TickTrain(s)
 if !s.FoundFacts[7] {t.Fatal("arrival effect not applied")}
 if !s.ProcessedStations["one:arrival"] {t.Fatal("arrival not marked processed")}
 TickTrain(s)
 if len(s.ProcessedStations)!=1 {t.Fatal("duplicate station processing")}
}

func TestTimedEventExpiresToDefaultChoice(t *testing.T){
 s:=&GameState{Clock:NewClock(time.Now()),FoundFacts:make(map[int]bool),
 Events:[]Event{{ID:1,TimeLimit:10,DefaultOptionID:2,Options:[]EventOption{
  {ID:2,Text:"Do nothing",Action:Action{Effects:[]Effect{{Type:"discover_fact",ID:5}}}},
 }}}}
 if !TriggerEvent(s,1){t.Fatal("failed to trigger")}
 if ExpireActiveEvent(s){t.Fatal("event expired too early")}
 s.EventDeadline=time.Now().Add(-time.Second)
 if !ExpireActiveEvent(s){t.Fatal("event not expired")}
 if !s.FoundFacts[5] || s.ActiveEvent!=nil {t.Fatal("default outcome missing")}
}
func TestScheduledEventStartsOnce(t *testing.T){
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 s:=&GameState{Clock:NewClock(start),RouteStart:start,
 Events:[]Event{{ID:4,At:"00:01",TimeLimit:10,Options:[]EventOption{{ID:1,Text:"Wait"}}}}}
 TickEvents(s)
 if s.ActiveEvent!=nil {t.Fatal("event too early")}
 s.Clock.Add(12*time.Minute)
 TickEvents(s)
 if s.ActiveEvent==nil || *s.ActiveEvent!=4 {t.Fatal("scheduled event missing")}
 if !ChooseEventOption(s,1){t.Fatal("option rejected")}
 TickEvents(s)
 if s.ActiveEvent!=nil {t.Fatal("event triggered twice")}
}

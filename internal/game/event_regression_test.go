package game

import (
 "testing"
 "time"
)

func TestEventChoiceCanTriggerFollowUpEvent(t *testing.T){
 s:=&GameState{Events:[]Event{
  {ID:1,Options:[]EventOption{{ID:10,Action:Action{Effects:[]Effect{{Type:"trigger_event",ID:2}}}}}},
  {ID:2,TimeLimit:60,Options:[]EventOption{{ID:20}}},
 }}
 if !TriggerEvent(s,1){t.Fatal("first event not triggered")}
 if !ChooseEventOption(s,10){t.Fatal("choice failed")}
 if s.ActiveEvent==nil || *s.ActiveEvent!=2 {t.Fatalf("follow-up event lost: %v",s.ActiveEvent)}
 if s.EventDeadline.IsZero(){t.Fatal("follow-up deadline was erased")}
}

func TestExpiredEventCanTriggerFollowUpEvent(t *testing.T){
 s:=&GameState{Events:[]Event{
  {ID:1,TimeLimit:1,DefaultOptionID:10,Options:[]EventOption{{ID:10,Action:Action{Effects:[]Effect{{Type:"trigger_event",ID:2}}}}}},
  {ID:2,TimeLimit:60,Options:[]EventOption{{ID:20}}},
 }}
 if !TriggerEvent(s,1){t.Fatal("first event not triggered")}
 s.EventDeadline=time.Now().Add(-time.Second)
 if !ExpireActiveEvent(s){t.Fatal("event did not expire")}
 if s.ActiveEvent==nil || *s.ActiveEvent!=2{t.Fatal("default consequence lost follow-up event")}
}

func TestExpiredEventRejectsLateChoiceAndAppliesDefault(t *testing.T){
 s:=&GameState{
  FoundFacts:map[int]bool{},
  Events:[]Event{{ID:1,TimeLimit:5,DefaultOptionID:20,Options:[]EventOption{
   {ID:10,Action:Action{Effects:[]Effect{{Type:"discover_fact",ID:10}}}},
   {ID:20,Action:Action{Effects:[]Effect{{Type:"discover_fact",ID:20}}}},
  }}},
 }
 if !TriggerEvent(s,1){t.Fatal("trigger failed")}
 s.EventDeadline=time.Now().Add(-time.Second)
 if ChooseEventOption(s,10){t.Fatal("late choice accepted")}
 if s.FoundFacts[10] || !s.FoundFacts[20]{t.Fatal("wrong consequence after deadline")}
 if s.ActiveEvent!=nil{t.Fatal("expired event still active")}
}

func TestPausedClockDoesNotAdvanceWithAddedTime(t *testing.T){
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 clock:=NewClock(start)
 defer clock.StopTime()
 clock.Pause()
 clock.Add(45*time.Second)
 if got:=clock.GetTime();!got.Equal(start.Add(45*time.Second)){
  t.Fatalf("explicit action time cost lost during pause: %v",got)
 }
 if !clock.IsPaused(){t.Fatal("explicit time addition resumed paused clock")}
}

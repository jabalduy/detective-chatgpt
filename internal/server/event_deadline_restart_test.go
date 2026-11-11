package server

import (
 "detective/internal/game"
 "testing"
 "time"
)

func TestUrgentEventDeadlineSurvivesRestore(t *testing.T) {
 now:=time.Now()
 start:=time.Date(2026,1,1,12,0,0,0,time.UTC)
 id:=1
 original:=&game.GameState{
  Case:game.CaseInfo{ID:"test"},
  Clock:game.NewClock(start),
  ActiveEvent:&id,
  EventDeadline:now.Add(-time.Second),
  FoundFacts:map[int]bool{},
  Events:[]game.Event{{ID:1,TimeLimit:10,DefaultOptionID:20,Options:[]game.EventOption{
   {ID:10,Action:game.Action{Effects:[]game.Effect{{Type:"discover_fact",ID:10}}}},
   {ID:20,Action:game.Action{Effects:[]game.Effect{{Type:"discover_fact",ID:20}}}},
  }}},
 }
 defer original.Clock.StopTime()
 saved:=snapshot(original)
 restored:=&game.GameState{Clock:game.NewClock(start),Events:original.Events}
 defer restored.Clock.StopTime()
 applySave(restored,saved)
 if !restored.EventDeadline.Equal(original.EventDeadline) {t.Fatal("restore reset wall-clock deadline")}
 if game.ChooseEventOption(restored,10) {t.Fatal("late decision accepted after restart")}
 if !restored.FoundFacts[20] || restored.FoundFacts[10] {t.Fatal("default consequence not applied")}
 if restored.ActiveEvent!=nil {t.Fatal("expired event remains active")}
}

func TestUrgentEventRemainingTimeNotRefilled(t *testing.T) {
 start:=time.Date(2026,1,1,12,0,0,0,time.UTC)
 id:=1
 deadline:=time.Now().Add(2*time.Second)
 original:=&game.GameState{Case:game.CaseInfo{ID:"test"},Clock:game.NewClock(start),ActiveEvent:&id,EventDeadline:deadline}
 defer original.Clock.StopTime()
 saved:=snapshot(original)
 restored:=&game.GameState{Clock:game.NewClock(start)}
 defer restored.Clock.StopTime()
 applySave(restored,saved)
 if !restored.EventDeadline.Equal(deadline) {t.Fatalf("deadline moved from %v to %v",deadline,restored.EventDeadline)}
}

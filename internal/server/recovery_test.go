package server

import (
 "detective/internal/game"
 "testing"
 "time"
)

func TestUrgentEventDeadlineSurvivesRestart(t *testing.T){
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 id:=42
 original:=&game.GameState{
  Case:game.CaseInfo{ID:"test-case"},Clock:game.NewClock(start),
  ActiveEvent:&id,EventDeadline:time.Now().Add(30*time.Second),
 }
 saved:=snapshot(original)
 original.Clock.StopTime()
 if saved.EventRemainingSeconds<29 || saved.EventRemainingSeconds>30{
  t.Fatalf("remaining time not captured: %v",saved.EventRemainingSeconds)
 }
 // Simulate a long server outage: the original wall-clock deadline expired.
 saved.EventDeadline=time.Now().Add(-time.Hour)
 restored:=&game.GameState{Clock:game.NewClock(start)}
 applySave(restored,saved)
 defer restored.Clock.StopTime()
 if restored.ActiveEvent==nil || *restored.ActiveEvent!=id{t.Fatal("active event lost")}
 remaining:=time.Until(restored.EventDeadline)
 if remaining>=0 || remaining< -time.Hour-time.Second {
  t.Fatalf("expired deadline should not be reset after downtime: %v",remaining)
 }
}

func TestExpiredUrgentEventDoesNotGainExtraTime(t *testing.T){
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 id:=42
 original:=&game.GameState{Case:game.CaseInfo{ID:"test-case"},Clock:game.NewClock(start),ActiveEvent:&id,EventDeadline:time.Now().Add(-time.Second)}
 saved:=snapshot(original)
 original.Clock.StopTime()
 restored:=&game.GameState{Clock:game.NewClock(start)}
 applySave(restored,saved)
 defer restored.Clock.StopTime()
 if time.Until(restored.EventDeadline)>time.Second{t.Fatal("expired event received extra time")}
}

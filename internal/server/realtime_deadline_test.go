package server

import (
 "testing"
 "time"
 "detective/internal/game"
)

func TestRealTimeDeadlinesSurviveSaveRestore(t *testing.T){
 start:=time.Date(2026,11,9,12,0,0,0,time.UTC)
 s:=&game.GameState{
  Case:game.CaseInfo{ID:"test"},
  Clock:game.NewClock(start),
  InvestigationDeadline:start.Add(60*time.Minute),
  AccusationDeadline:start.Add(65*time.Minute),
 }
 defer s.Clock.StopTime()
 saved:=snapshot(s)
 if !saved.InvestigationDeadline.Equal(s.InvestigationDeadline) || !saved.AccusationDeadline.Equal(s.AccusationDeadline){
  t.Fatal("snapshot lost real-time deadlines")
 }
 restored:=&game.GameState{Clock:game.NewClock(start)}
 defer restored.Clock.StopTime()
 applySave(restored,saved)
 if !restored.InvestigationDeadline.Equal(s.InvestigationDeadline) || !restored.AccusationDeadline.Equal(s.AccusationDeadline){
  t.Fatal("restore changed absolute deadlines")
 }
}

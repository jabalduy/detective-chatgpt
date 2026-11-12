package game

import (
 "testing"
 "time"
)

func TestPrologueCannotUseInvestigationOrAccusationActions(t *testing.T) {
 now:=time.Date(2026,11,12,12,0,0,0,time.UTC)
 s:=&GameState{NarrativeAct:ActPrologue}
 if s.CanInvestigateAt(now) || s.CanAccuseAt(now) {
  t.Fatal("prologue leaked investigation or accusation actions")
 }
 // Old sessions with an empty narrative act must keep their prior behavior.
 s.NarrativeAct=""
 if !s.CanInvestigateAt(now) {
  t.Fatal("legacy investigation became inaccessible")
 }
 s.NarrativeAct=ActInvestigation
 s.InvestigationDeadline=now.Add(-time.Minute)
 s.AccusationDeadline=now.Add(time.Minute)
 if s.CanInvestigateAt(now) || !s.CanAccuseAt(now) {
  t.Fatal("accusation window behavior regressed")
 }
 s.NarrativeAct=ActPrologue
 if s.CanAccuseAt(now) {t.Fatal("prologue leaked accusation window")}
}

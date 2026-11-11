package game

import (
 "testing"
 "time"
)

func TestRealTimePhases(t *testing.T) {
 start:=time.Date(2026,11,9,12,0,0,0,time.UTC)
 s:=&GameState{InvestigationDeadline:start.Add(60*time.Minute),AccusationDeadline:start.Add(65*time.Minute)}
 cases:=[]struct{at time.Time; want RealTimePhase}{
  {start,PhaseInvestigation},
  {s.InvestigationDeadline.Add(-time.Nanosecond),PhaseInvestigation},
  {s.InvestigationDeadline,PhaseAccusation},
  {s.AccusationDeadline.Add(-time.Nanosecond),PhaseAccusation},
  {s.AccusationDeadline,PhaseExpired},
 }
 for _,tc:=range cases {if got:=s.PhaseAt(tc.at);got!=tc.want{t.Errorf("at %v got %s want %s",tc.at,got,tc.want)}}
 s.ActiveConversation="witness"
 if got:=s.PhaseAt(s.InvestigationDeadline);got!=PhaseConversationGrace{t.Fatalf("conversation grace: %s",got)}
 if s.CanInvestigateAt(s.InvestigationDeadline){t.Fatal("new investigation action allowed in grace")}
 if s.CanAccuseAt(s.InvestigationDeadline){t.Fatal("accusation allowed while conversation active")}
 if got:=s.PhaseAt(s.AccusationDeadline);got!=PhaseConversationGrace{t.Fatalf("active conversation should retain grace: %s",got)}
 s.ActiveConversation=""
 if !s.CanAccuseAt(s.InvestigationDeadline){t.Fatal("accusation not enabled after conversation")}
 s.Ending=&GameEnding{}
 if got:=s.PhaseAt(start);got!=PhaseFinished{t.Fatalf("ending phase: %s",got)}
}

package game

import (
 "errors"
 "testing"
 "time"
)

func TestBeginInvestigationOnce(t *testing.T) {
 clock:=NewClock(time.Date(1928,10,1,22,16,59,0,time.UTC))
 defer clock.StopTime()
 s:=&GameState{NarrativeAct:ActPrologue,Clock:clock}
 now:=time.Date(2026,11,12,10,0,0,0,time.UTC)
 if changed,err:=s.BeginInvestigation(now);changed || !errors.Is(err,ErrDiscoveryTooEarly) {
  t.Fatalf("early transition: changed=%v err=%v",changed,err)
 }
 if !s.InvestigationDeadline.IsZero() {t.Fatal("early call armed deadline")}
 clock.Add(time.Second)
 changed,err:=s.BeginInvestigation(now)
 if err!=nil || !changed {t.Fatalf("discovery failed: changed=%v err=%v",changed,err)}
 if s.CurrentNarrativeAct()!=ActInvestigation || !s.DiscoveryCompleted {t.Fatal("act not advanced")}
 deadline:=now.Add(time.Hour)
 if !s.InvestigationDeadline.Equal(deadline) || !s.AccusationDeadline.Equal(deadline.Add(5*time.Minute)) {
  t.Fatal("incorrect investigation/accusation deadlines")
 }
 changed,err=s.BeginInvestigation(now.Add(10*time.Minute))
 if err!=nil || changed || !s.InvestigationDeadline.Equal(deadline) {
  t.Fatal("duplicate transition extended deadlines")
 }
}

func TestBeginInvestigationWaitsForUrgentChoice(t *testing.T) {
 clock:=NewClock(time.Date(1928,10,1,22,17,0,0,time.UTC))
 defer clock.StopTime()
 id:=105
 s:=&GameState{NarrativeAct:ActPrologue,Clock:clock,ActiveEvent:&id}
 if changed,err:=s.BeginInvestigation(time.Now());changed || !errors.Is(err,ErrDiscoveryInteractionActive) {
  t.Fatalf("active chase interrupted: %v %v",changed,err)
 }
 s.ActiveEvent=nil
 if changed,err:=s.BeginInvestigation(time.Now());err!=nil || !changed {
  t.Fatalf("transition after chase: %v %v",changed,err)
 }
}

func TestLegacyInvestigationDoesNotRestart(t *testing.T) {
 s:=&GameState{}
 if changed,err:=s.BeginInvestigation(time.Now());err!=nil || changed {
  t.Fatalf("legacy investigation unexpectedly transitioned: %v %v",changed,err)
 }
}

func TestBeginInvestigationReplacesStalePrologueDeadlines(t *testing.T) {
 clock:=NewClock(time.Date(1928,10,1,22,17,0,0,time.UTC))
 defer clock.StopTime()
 now:=time.Date(2026,11,12,12,0,0,0,time.UTC)
 s:=&GameState{
  NarrativeAct:ActPrologue,Clock:clock,
  InvestigationDeadline:now.Add(-time.Hour),
  AccusationDeadline:now.Add(-55*time.Minute),
 }
 changed,err:=s.BeginInvestigation(now)
 if err!=nil || !changed {t.Fatalf("transition failed: changed=%v err=%v",changed,err)}
 want:=now.Add(time.Hour)
 if !s.InvestigationDeadline.Equal(want) {t.Fatalf("deadline=%v want %v",s.InvestigationDeadline,want)}
 if !s.AccusationDeadline.Equal(want.Add(5*time.Minute)) {t.Fatalf("accusation deadline=%v",s.AccusationDeadline)}
 changed,err=s.BeginInvestigation(now.Add(10*time.Minute))
 if err!=nil || changed || !s.InvestigationDeadline.Equal(want) {t.Fatal("repeat transition modified deadline")}
}

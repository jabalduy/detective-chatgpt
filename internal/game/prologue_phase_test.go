package game

import (
 "testing"
 "time"
)

func TestPhaseAtDistinguishesPrologueFromInvestigation(t *testing.T) {
 now:=time.Now()
 for _,tc:=range []struct{
  name string
  state GameState
  want RealTimePhase
 }{
  {"prologue without deadline",GameState{NarrativeAct:ActPrologue},PhasePrologue},
  {"prologue with expired deadline",GameState{NarrativeAct:ActPrologue,InvestigationDeadline:now.Add(-time.Minute)},PhasePrologue},
  {"investigation without deadline",GameState{NarrativeAct:ActInvestigation},PhaseInvestigation},
  {"legacy without deadline",GameState{},PhaseInvestigation},
  {"finished takes priority",GameState{NarrativeAct:ActPrologue,Ending:&GameEnding{}},PhaseFinished},
 } {
  t.Run(tc.name,func(t *testing.T){
   if got:=tc.state.PhaseAt(now);got!=tc.want {t.Fatalf("phase=%s want %s",got,tc.want)}
  })
 }
}

package server

import (
 "path/filepath"
 "testing"
 "time"

 "detective/internal/game"
)

func TestNarrativeActSaveRestoreAndLegacyFallback(t *testing.T) {
 start:=time.Date(1928,10,1,21,0,0,0,time.UTC)
 original:=&game.GameState{
  Case:game.CaseInfo{ID:"case_express_1928"},
  Clock:game.NewClock(start),
  NarrativeAct:game.ActPrologue,
  DiscoveryCompleted:false,
  CurrentScene:"corridor",
 }
 defer original.Clock.StopTime()
 path:=filepath.Join(t.TempDir(),"prologue.json")
 if err:=writeSave(path,snapshot(original));err!=nil {t.Fatal(err)}
 saved,err:=readSave(path)
 if err!=nil {t.Fatal(err)}
 restored:=&game.GameState{Clock:game.NewClock(start)}
 defer restored.Clock.StopTime()
 applySave(restored,saved)
 if restored.CurrentNarrativeAct()!=game.ActPrologue || restored.DiscoveryCompleted {
  t.Fatal("prologue state lost after restore")
 }
 if !restored.InvestigationDeadline.IsZero() {t.Fatal("prologue unexpectedly started investigation deadline")}
 original.NarrativeAct=game.ActInvestigation
 original.DiscoveryCompleted=true
 original.InvestigationDeadline=time.Now().Add(time.Hour)
 saved=snapshot(original)
 applySave(restored,saved)
 if restored.CurrentNarrativeAct()!=game.ActInvestigation || !restored.DiscoveryCompleted {
  t.Fatal("discovery state lost")
 }
 if !restored.InvestigationDeadline.Equal(original.InvestigationDeadline) {
  t.Fatal("investigation deadline changed on restore")
 }
 legacy:=saveData{Version:saveVersion,CaseID:"case_express_1928",ClockTime:start}
 if err:=writeSave(path,legacy);err!=nil {t.Fatal(err)}
 legacy,err=readSave(path)
 if err!=nil {t.Fatal(err)}
 applySave(restored,legacy)
 if restored.CurrentNarrativeAct()!=game.ActInvestigation || restored.DiscoveryCompleted {
  t.Fatal("legacy save must remain in investigation, without fake discovery")
 }
}

func TestNarrativeActFinishedOverridesStoryState(t *testing.T) {
 s:=&game.GameState{NarrativeAct:game.ActPrologue,Ending:&game.GameEnding{}}
 if got:=s.CurrentNarrativeAct();got!=game.ActFinished {
  t.Fatalf("finished game act=%q",got)
 }
}

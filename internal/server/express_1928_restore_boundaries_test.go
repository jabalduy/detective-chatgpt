package server

import (
 "path/filepath"
 "testing"
 "time"

 "detective/internal/cases"
 "detective/internal/game"
)

// An active ten-second event must not gain time when a save is restored.
func TestExpress1928RestoreDoesNotExtendUrgentDecision(t *testing.T) {
 def,err:=cases.LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 start:=time.Date(1928,10,1,21,0,0,0,time.UTC)
 deadline:=time.Now().Add(-time.Second)
 original:=&game.GameState{
  Case:def.CaseInfo,Clock:game.NewClock(start),Scenes:def.Scenes,
  Events:def.Events,CurrentScene:"corridor",
  TriggeredEvents:map[int]bool{101:true},
  ActiveEvent:new(int),EventDeadline:deadline,
  FoundFacts:map[int]bool{},
 }
 *original.ActiveEvent=101
 defer original.Clock.StopTime()
 path:=filepath.Join(t.TempDir(),"urgent.json")
 if err:=writeSave(path,snapshot(original));err!=nil {t.Fatal(err)}
 saved,err:=readSave(path);if err!=nil {t.Fatal(err)}
 restored:=&game.GameState{
  Case:def.CaseInfo,Clock:game.NewClock(start),Events:def.Events,
 }
 defer restored.Clock.StopTime()
 applySave(restored,saved)
 if restored.ActiveEvent==nil || *restored.ActiveEvent!=101 {t.Fatal("urgent event lost")}
 if !restored.EventDeadline.Equal(deadline) {t.Fatal("event deadline was reset")}
 if time.Until(restored.EventDeadline)>0 {t.Fatal("expired event gained new time")}
 if !game.ExpireActiveEvent(restored) {t.Fatal("expired event cannot resolve")}
 if restored.ActiveEvent!=nil {t.Fatal("expired event remained active")}
 if restored.FoundFacts[17] {t.Fatal("default action awarded optional evidence")}
}

// Reloading a partially solved case must preserve evidence, the selected
// room, and the ability to finish deductions after the restart.
func TestExpress1928RestorePartialDeductionAndContinue(t *testing.T) {
 def,err:=cases.LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 start:=time.Date(1928,10,1,21,0,0,0,time.UTC)
 original:=&game.GameState{
  Case:def.CaseInfo,Clock:game.NewClock(start),Scenes:def.Scenes,
  Deductions:def.Deductions,CurrentScene:"archive",
  FoundFacts:map[int]bool{1:true,8:true,2:true,3:true,6:true,7:true},
  FoundMotives:map[int]bool{},SolvedDeductions:map[int]bool{},
  InvestigationDeadline:time.Now().Add(20*time.Minute),
  AccusationDeadline:time.Now().Add(25*time.Minute),
 }
 defer original.Clock.StopTime()
 for _,d:=range original.Deductions {
  if d.ID==1 && !game.SolveDeduction(original,d) {t.Fatal("first deduction failed")}
 }
 path:=filepath.Join(t.TempDir(),"partial.json")
 if err:=writeSave(path,snapshot(original));err!=nil {t.Fatal(err)}
 saved,err:=readSave(path);if err!=nil {t.Fatal(err)}
 restored:=&game.GameState{
  Case:def.CaseInfo,Clock:game.NewClock(start),
  Scenes:def.Scenes,Deductions:def.Deductions,
 }
 defer restored.Clock.StopTime()
 applySave(restored,saved)
 if restored.CurrentScene!="archive" || !restored.SolvedDeductions[1] || !restored.FoundFacts[10] {
  t.Fatal("restored partial investigation lost its progress")
 }
 if !restored.InvestigationDeadline.Equal(original.InvestigationDeadline) ||
  !restored.AccusationDeadline.Equal(original.AccusationDeadline) {
  t.Fatal("restore extended investigation deadlines")
 }
 for _,id:=range []int{2,4} {
  found:=false
  for _,d:=range restored.Deductions {
   if d.ID!=id {continue}
   found=true
   if !game.SolveDeduction(restored,d) {t.Fatalf("deduction %d failed after restore",id)}
  }
  if !found {t.Fatalf("deduction %d missing",id)}
 }
 if !restored.FoundFacts[11] || !restored.FoundMotives[1] {
  t.Fatal("restored investigation cannot continue")
 }
 for _,d:=range restored.Deductions {
  if d.ID==1 && game.SolveDeduction(restored,d) {
   t.Fatal("already completed deduction repeated after restore")
  }
 }
}

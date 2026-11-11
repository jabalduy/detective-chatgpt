package server

import (
 "path/filepath"
 "testing"
 "time"

 "detective/internal/cases"
 "detective/internal/game"
)

func TestExpress1928AutosaveRoundTripPreservesInvestigation(t *testing.T) {
 def,err:=cases.LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 start:=time.Date(1928,10,1,21,0,0,0,time.UTC)
 original:=&game.GameState{
  Case:def.CaseInfo,Clock:game.NewClock(start),
  Scenes:def.Scenes,Conversations:def.Conversations,Characters:def.Characters,
  Deductions:def.Deductions,Motives:def.Motives,EndingRules:def.EndingRules,
  CurrentScene:"corridor",FoundFacts:map[int]bool{},FoundMotives:map[int]bool{},
  SolvedDeductions:map[int]bool{},
  InvestigationDeadline:time.Now().Add(60*time.Minute),
  AccusationDeadline:time.Now().Add(65*time.Minute),
 }
 defer original.Clock.StopTime()
 for _,id:=range []string{"victor","ribbon","paperweight","watch","note","back","staff","box","nikola","cancellation"} {
  if err:=game.ChooseSceneOption(original,id);err!=nil {t.Fatalf("choice %s: %v",id,err)}
 }
 if err:=game.StartConversation(original,"marta_intro");err!=nil {t.Fatal(err)}
 original.Clock.Add(9*time.Minute)
 saved:=snapshot(original)
 path:=filepath.Join(t.TempDir(),"express.json")
 if err:=writeSave(path,saved);err!=nil {t.Fatal(err)}
 loaded,err:=readSave(path)
 if err!=nil {t.Fatal(err)}
 restored:=&game.GameState{
  Case:def.CaseInfo,Clock:game.NewClock(start),
  Scenes:def.Scenes,Conversations:def.Conversations,Characters:def.Characters,
  Deductions:def.Deductions,Motives:def.Motives,EndingRules:def.EndingRules,
 }
 defer restored.Clock.StopTime()
 applySave(restored,loaded)
 if restored.CurrentScene!="staff" || restored.ActiveConversation!="marta_intro" || !restored.Clock.IsPaused() {
  t.Fatal("scene, active dialogue or paused clock lost")
 }
 if !restored.Clock.GetTime().Equal(saved.ClockTime) {t.Fatal("route clock time changed on restore")}
 if !restored.InvestigationDeadline.Equal(original.InvestigationDeadline) || !restored.AccusationDeadline.Equal(original.AccusationDeadline) {
  t.Fatal("wall-clock deadlines were reset")
 }
 for _,id:=range []int{1,2,3,5,8,13,14} {
  if !restored.FoundFacts[id] {t.Errorf("fact %d lost on disk roundtrip",id)}
 }
 if _,err:=game.ChooseConversationOption(restored,"leave");err!=nil {t.Fatalf("cannot resume restored dialogue: %v",err)}
 if restored.Clock.IsPaused() {t.Fatal("clock did not resume after restored dialogue")}
 if err:=game.ChooseSceneOption(restored,"back");err!=nil {t.Fatal(err)}
 if restored.CurrentScene!="corridor" {t.Fatal("cannot continue investigation after restore")}
}

func TestExpress1928FinalVerdictPersistsAcrossRestart(t *testing.T){
 def,err:=cases.LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 start:=time.Date(1928,10,1,21,0,0,0,time.UTC)
 original:=&game.GameState{
  Case:def.CaseInfo,Clock:game.NewClock(start),Motives:def.Motives,
  FoundFacts:map[int]bool{10:true,11:true,12:true},
  FoundMotives:map[int]bool{1:true},
 }
 defer original.Clock.StopTime()
 ending,err:=game.SubmitAccusation(original,game.Accusation{SuspectID:1,MotiveID:1,FactIDs:[]int{10,11,12}})
 if err!=nil || ending.Type!="win" {t.Fatalf("verdict %+v: %v",ending,err)}
 path:=filepath.Join(t.TempDir(),"verdict.json")
 if err:=writeSave(path,snapshot(original));err!=nil {t.Fatal(err)}
 loaded,err:=readSave(path)
 if err!=nil {t.Fatal(err)}
 restored:=&game.GameState{Case:def.CaseInfo,Clock:game.NewClock(start),Motives:def.Motives}
 defer restored.Clock.StopTime()
 applySave(restored,loaded)
 if restored.Ending==nil || restored.Ending.Type!="win" || !restored.Clock.IsPaused() {
  t.Fatal("completed verdict or stopped clock lost after restart")
 }
 if _,err:=game.SubmitAccusation(restored,game.Accusation{SuspectID:1,MotiveID:1,FactIDs:[]int{10,11,12}});err==nil {
  t.Fatal("completed case allowed another accusation")
 }
}

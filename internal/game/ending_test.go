package game

import (
 "testing"
 "time"
)

func endingTestState()*GameState{
 return &GameState{
  Clock:NewClock(time.Date(2026,1,1,22,0,0,0,time.UTC)),
  Case:CaseInfo{Solution:Solution{SuspectID:1,MotiveID:2,FactIDs:[]int{3,4}},Endings:Endings{Win:"Solved",Fail:"Wrong",Unsolved:"Incomplete"}},
  Motives:[]Motive{{ID:2,SuspectID:1},{ID:5,SuspectID:1},{ID:6,SuspectID:9}},
  FoundMotives:map[int]bool{2:true,5:true,6:true},
  FoundFacts:map[int]bool{3:true,4:true},
 }
}
func TestAccusationVerdicts(t *testing.T){
 cases:=[]struct{name string;acc Accusation;want string}{
  {"win",Accusation{1,2,[]int{3,4}},"win"},
  {"wrong suspect",Accusation{9,6,[]int{3,4}},"fail"},
  {"wrong motive",Accusation{1,5,[]int{3,4}},"unsolved"},
  {"missing evidence",Accusation{1,2,[]int{3}},"unsolved"},
 }
 for _,tc:=range cases {t.Run(tc.name,func(t *testing.T){
  s:=endingTestState()
  got,err:=SubmitAccusation(s,tc.acc)
  if err!=nil {t.Fatal(err)}
  if got.Type!=tc.want {t.Fatalf("got %s want %s",got.Type,tc.want)}
  if s.Ending==nil {t.Fatal("ending not stored")}
  if _,err:=SubmitAccusation(s,tc.acc);err==nil {t.Fatal("second accusation accepted")}
  if PerformAction(s,Action{}) {t.Fatal("action after ending accepted")}
 })}
}
func TestAccusationRejectsUndiscoveredEvidence(t *testing.T){
 s:=endingTestState()
 _,err:=SubmitAccusation(s,Accusation{1,2,[]int{3,999}})
 if err==nil || s.Ending!=nil {t.Fatal("invalid accusation ended game")}
}
func TestConditionalEndingVariant(t *testing.T){
 s:=endingTestState()
 s.EndingRules=[]EndingRule{{ID:"win",Title:"Special",Text:"Special ending",Priority:10,Requirements:[]Requirement{{Type:"fact_found",ID:4}}}}
 got,err:=SubmitAccusation(s,Accusation{1,2,[]int{3,4}})
 if err!=nil || got.Title!="Special" || got.Text!="Special ending" {t.Fatalf("ending=%+v err=%v",got,err)}
}

func TestAccusationRejectsDuplicateEvidenceWithoutEnding(t *testing.T) {
 s:=endingTestState()
 defer s.Clock.StopTime()
 _,err:=SubmitAccusation(s,Accusation{SuspectID:1,MotiveID:2,FactIDs:[]int{3,3,4}})
 if err==nil || s.Ending!=nil {t.Fatal("duplicate evidence must not end the game")}
 if _,err:=SubmitAccusation(s,Accusation{SuspectID:1,MotiveID:2,FactIDs:[]int{4,3}});err!=nil {
  t.Fatalf("valid accusation after rejected attempt failed: %v",err)
 }
}

func TestEndingRulePriorityAndRequirements(t *testing.T) {
 s:=endingTestState()
 defer s.Clock.StopTime()
 s.EndingRules=[]EndingRule{
  {ID:"win",Title:"Low",Text:"Low priority",Priority:1},
  {ID:"win",Title:"Locked",Text:"Missing clue",Priority:100,Requirements:[]Requirement{{Type:"clue_found",ID:99}}},
  {ID:"win",Title:"High",Text:"High priority",Priority:5,Requirements:[]Requirement{{Type:"fact_found",ID:4}}},
  {ID:"fail",Title:"Wrong verdict",Priority:1000},
 }
 got,err:=SubmitAccusation(s,Accusation{SuspectID:1,MotiveID:2,FactIDs:[]int{3,4}})
 if err!=nil || got.Type!="win" || got.Title!="High" || got.Text!="High priority" {
  t.Fatalf("incorrect ending rule selection: %+v %v",got,err)
 }
}

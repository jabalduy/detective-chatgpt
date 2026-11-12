package cases

import (
 "path/filepath"
 "strings"
 "testing"

 "detective/internal/game"
)

func TestExpress1928FalseAccusationsHaveDistinctExplanations(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 tests:=[]struct{suspect,motive int;expected string}{
  {2,2,"Борис"},
  {5,3,"Клара"},
  {6,4,"Макс"},
  {4,5,"Граф"},
 }
 seen:=map[string]bool{}
 for _,tc:=range tests {
  s:=&game.GameState{
   Case:def.CaseInfo,Motives:def.Motives,EndingRules:def.EndingRules,
   FoundMotives:map[int]bool{tc.motive:true},
   FoundFacts:map[int]bool{},
  }
  ending,err:=game.SubmitAccusation(s,game.Accusation{SuspectID:tc.suspect,MotiveID:tc.motive})
  if err!=nil {t.Fatal(err)}
  if ending.Type!="fail" {t.Fatalf("suspect %d verdict %q",tc.suspect,ending.Type)}
  if !strings.Contains(ending.Title,tc.expected) {t.Errorf("suspect %d generic ending: %q",tc.suspect,ending.Title)}
  if seen[ending.Text] {t.Errorf("suspect %d shares another ending",tc.suspect)}
  seen[ending.Text]=true
 }
}

func TestEndingRuleMatchesAccusedSuspectNotOtherEvidence(t *testing.T) {
 s:=&game.GameState{
  Case:game.CaseInfo{},
  Motives:[]game.Motive{{ID:1,SuspectID:2},{ID:2,SuspectID:3}},
  FoundMotives:map[int]bool{1:true,2:true},
  FoundFacts:map[int]bool{},
  EndingRules:[]game.EndingRule{
   {ID:"fail",Title:"generic",Text:"generic",Priority:1},
   {ID:"fail",SuspectID:2,Title:"Boris",Text:"Boris ending",Priority:20},
  },
 }
 ending,err:=game.SubmitAccusation(s,game.Accusation{SuspectID:3,MotiveID:2})
 if err!=nil {t.Fatal(err)}
 if ending.Title!="generic" {t.Fatalf("other suspect received Boris ending: %q",ending.Title)}
}

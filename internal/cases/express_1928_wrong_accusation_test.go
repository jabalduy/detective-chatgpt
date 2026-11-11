package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

func TestExpress1928WrongAccusationIsReachable(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,issue:=range ValidateCaseIntegrity(def) {t.Error(issue)}
 s:=&game.GameState{
  Case:def.CaseInfo,Motives:def.Motives,Deductions:def.Deductions,
  EndingRules:def.EndingRules,FoundFacts:map[int]bool{4:true,13:true},
  FoundMotives:map[int]bool{},SolvedDeductions:map[int]bool{},
 }
 var boris *game.Deduction
 for i:=range s.Deductions {if s.Deductions[i].ID==7 {boris=&s.Deductions[i]}}
 if boris==nil {t.Fatal("Boris suspicion deduction missing")}
 if !game.SolveDeduction(s,*boris) {t.Fatal("cannot deduce Boris suspicion")}
 if !s.FoundMotives[2] {t.Fatal("Boris motive not discovered")}
 ending,err:=game.SubmitAccusation(s,game.Accusation{SuspectID:2,MotiveID:2,FactIDs:[]int{4,13}})
 if err!=nil {t.Fatal(err)}
 if ending.Type!="fail" {t.Fatalf("wrong suspect verdict = %q, want fail",ending.Type)}
 if ending.Title=="" || ending.Text=="" {t.Fatal("wrong accusation narrative missing")}
}

func TestExpress1928BorisSuspicionRequiresBothFacts(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,facts:=range []map[int]bool{{},{4:true},{13:true}} {
  s:=&game.GameState{FoundFacts:facts,FoundMotives:map[int]bool{},SolvedDeductions:map[int]bool{}}
  for _,d:=range def.Deductions {
   if d.ID==7 && game.SolveDeduction(s,d) {t.Fatalf("Boris suspicion unlocked prematurely with %v",facts)}
  }
 }
}

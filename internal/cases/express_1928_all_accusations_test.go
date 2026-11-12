package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

// Every authored false motive must lead to a legitimate losing verdict, not
// a request-validation error. The canonical murderer remains unchanged.
func TestExpress1928EveryFalseSuspectCanBeAccused(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 cases:=[]struct{suspect,motive int;facts []int;deduction int}{
  {2,2,[]int{4,13},7},
  {5,3,[]int{6,20},8},
  {6,4,[]int{21,22},9},
  {4,5,[]int{6,23},10},
 }
 for _,tc:=range cases {
  t.Run(string(rune('0'+tc.suspect)),func(t *testing.T){
   found:=map[int]bool{}
   for _,id:=range tc.facts {found[id]=true}
   s:=&game.GameState{
    Case:def.CaseInfo,Motives:def.Motives,Deductions:def.Deductions,
    EndingRules:def.EndingRules,FoundFacts:found,
    FoundMotives:map[int]bool{},SolvedDeductions:map[int]bool{},
   }
   var d *game.Deduction
   for i:=range s.Deductions {if s.Deductions[i].ID==tc.deduction {d=&s.Deductions[i];break}}
   if d==nil {t.Fatal("missing deduction")}
   if !game.SolveDeduction(s,*d) {t.Fatal("false motive cannot be deduced")}
   if !s.FoundMotives[tc.motive] {t.Fatal("false motive not awarded")}
   ending,err:=game.SubmitAccusation(s,game.Accusation{
    SuspectID:tc.suspect,MotiveID:tc.motive,FactIDs:tc.facts,
   })
   if err!=nil {t.Fatal(err)}
   if ending.Type!="fail" {t.Fatalf("verdict=%q want fail",ending.Type)}
   if ending.Title=="" || ending.Text=="" {t.Fatal("empty ending")}
   if _,err:=game.SubmitAccusation(s,game.Accusation{
    SuspectID:tc.suspect,MotiveID:tc.motive,FactIDs:tc.facts,
   });err==nil {t.Fatal("second accusation accepted after ending")}
  })
 }
}

func TestExpress1928AccusationRejectsUnprovenAndMismatchedMotives(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,tc:=range []struct{name string;suspect,motive int;facts []int;foundMotives map[int]bool}{
  {"unknown motive",5,3,[]int{6,20},map[int]bool{}},
  {"wrong suspect for known motive",5,2,[]int{4,13},map[int]bool{2:true}},
  {"missing evidence",5,3,[]int{6,20},map[int]bool{3:true}},
  {"duplicate evidence",5,3,[]int{6,6},map[int]bool{3:true}},
 } {
  t.Run(tc.name,func(t *testing.T){
   s:=&game.GameState{
    Case:def.CaseInfo,Motives:def.Motives,EndingRules:def.EndingRules,
    FoundMotives:tc.foundMotives,FoundFacts:map[int]bool{6:true,4:true,13:true},
   }
   _,err:=game.SubmitAccusation(s,game.Accusation{SuspectID:tc.suspect,MotiveID:tc.motive,FactIDs:tc.facts})
   if err==nil {t.Fatal("invalid accusation accepted")}
   if s.Ending!=nil {t.Fatal("invalid accusation ended investigation")}
  })
 }
}

func TestExpress1928CompleteAndIncompleteMartaAccusations(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,tc:=range []struct{name string;submitted []int;want string}{
  {"complete",[]int{10,11,12},"win"},
  {"missing clock",[]int{10,11},"unsolved"},
  {"missing weapon",[]int{10,12},"unsolved"},
  {"missing presence",[]int{11,12},"unsolved"},
 }{
  t.Run(tc.name,func(t *testing.T){
   s:=&game.GameState{
    Case:def.CaseInfo,Motives:def.Motives,EndingRules:def.EndingRules,
    FoundMotives:map[int]bool{1:true},
    FoundFacts:map[int]bool{10:true,11:true,12:true,16:true},
   }
   ending,err:=game.SubmitAccusation(s,game.Accusation{SuspectID:1,MotiveID:1,FactIDs:tc.submitted})
   if err!=nil {t.Fatal(err)}
   if ending.Type!=tc.want {t.Fatalf("verdict=%q want %q",ending.Type,tc.want)}
  })
 }
}

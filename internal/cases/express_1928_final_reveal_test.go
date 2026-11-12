package cases

import (
 "path/filepath"
 "strings"
 "testing"

 "detective/internal/game"
)

func TestExpress1928WinningRevealReconstructsCrime(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 s:=&game.GameState{
  Case:def.CaseInfo,Motives:def.Motives,EndingRules:def.EndingRules,
  FoundMotives:map[int]bool{1:true},
  FoundFacts:map[int]bool{10:true,11:true,12:true,16:true},
 }
 ending,err:=game.SubmitAccusation(s,game.Accusation{
  SuspectID:1,MotiveID:1,FactIDs:[]int{10,11,12},
 })
 if err!=nil {t.Fatal(err)}
 if ending.Type!="win" {t.Fatalf("verdict=%q",ending.Type)}
 for _,detail:=range []string{
  "21:28","21:35","21:38","21:39","21:41","22:03","22:06",
  "Павл","лента","Николя","пресс-папье","Борис","Софи","Марта",
 } {
  if !strings.Contains(ending.Text,detail) {t.Errorf("winning reveal missing %q",detail)}
 }
 if strings.Index(ending.Text,"21:28")>=strings.Index(ending.Text,"21:39") {
  t.Error("crime reconstruction is not chronological")
 }
}

func TestExpress1928IncompleteProofDoesNotSpoilFullReveal(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 s:=&game.GameState{
  Case:def.CaseInfo,Motives:def.Motives,EndingRules:def.EndingRules,
  FoundMotives:map[int]bool{1:true},
  FoundFacts:map[int]bool{10:true,11:true,12:true,16:true},
 }
 ending,err:=game.SubmitAccusation(s,game.Accusation{
  SuspectID:1,MotiveID:1,FactIDs:[]int{10,11},
 })
 if err!=nil {t.Fatal(err)}
 if ending.Type!="unsolved" {t.Fatalf("verdict=%q",ending.Type)}
 if strings.Contains(ending.Text,"21:39") {t.Fatal("incomplete accusation exposed full reconstruction")}
}

package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

// This is a data-driven playthrough, not merely a JSON schema test.
// All timed events are deliberately ignored: optional events must not gate the solution.
func TestExpress1928CanSolveWithoutTimedEvents(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 s:=&game.GameState{
  Scenes:def.Scenes,Conversations:def.Conversations,Characters:def.Characters,
  Deductions:def.Deductions,Events:def.Events,
  CurrentScene:"corridor",
  FoundFacts:map[int]bool{},FoundClues:map[int]bool{},
  SolvedDeductions:map[int]bool{},FoundMotives:map[int]bool{},
  TriggeredEvents:map[int]bool{},
  OpenDialogues:map[game.DialogueID]bool{},
 }
 move:=func(id string){
  t.Helper()
  if err:=game.ChooseSceneOption(s,id);err!=nil {t.Fatalf("scene %q choice %q: %v",s.CurrentScene,id,err)}
 }
 move("victor")
 for _,id:=range []string{"ribbon","paperweight","watch","note"}{move(id)}
 move("back");move("staff")
 for _,id:=range []string{"box","nikola","cancellation"}{move(id)}
 move("back");move("salon")
 for _,id:=range []string{"boris","rosa"}{move(id)}
 if err:=game.StartConversation(s,"boris_intro");err!=nil {t.Fatal(err)}
 if err:=game.ChooseConversationOption(s,"leave");err!=nil {t.Fatal(err)}
 move("back");move("archive")
 for _,id:=range []string{"fraud","family"}{move(id)}
 for _,id:=range []int{1,2,4,5}{
  var d *game.Deduction
  for i:=range s.Deductions {if s.Deductions[i].ID==id {d=&s.Deductions[i];break}}
  if d==nil || !game.SolveDeduction(s,*d) {t.Fatalf("deduction %d cannot be solved",id)}
 }
 if !s.FoundFacts[10] || !s.FoundFacts[11] || !s.FoundFacts[15] || !s.FoundMotives[1] {
  t.Fatal("missing intermediate solution evidence or motive")
 }
 move("back");move("salon")
 if err:=game.StartConversation(s,"boris_intro");err!=nil {t.Fatal(err)}
 view,err:=game.CurrentConversationView(s)
 if err!=nil {t.Fatal(err)}
 available:=false
 for _,c:=range view.Choices {if c.ID=="clock_confront" {available=true}}
 if !available {t.Fatal("Boris confession unavailable after required evidence")}
 if _,err:=game.ChooseConversationOption(s,"clock_confront");err!=nil {t.Fatal(err)}
 if _,err:=game.ChooseConversationOption(s,"record");err!=nil {t.Fatal(err)}
 var d *game.Deduction
 for i:=range s.Deductions {if s.Deductions[i].ID==6 {d=&s.Deductions[i];break}}
 if d==nil || !game.SolveDeduction(s,*d) {t.Fatal("final clock deduction cannot be solved")}
 for _,id:=range def.Solution.FactIDs {
  if !s.FoundFacts[id] {t.Errorf("required solution fact %d not reachable",id)}
 }
 if s.ActiveEvent!=nil {t.Fatal("an optional timed event was required")}
}

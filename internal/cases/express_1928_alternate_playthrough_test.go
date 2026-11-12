package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

// Explore the train out of the intended order, using alternate evidence
// locations and revisiting conversations. No timed event is required.
func TestExpress1928AlternateRouteStillReachesVerdict(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 s:=&game.GameState{
  Case:def.CaseInfo,Scenes:def.Scenes,Conversations:def.Conversations,
  Characters:def.Characters,Motives:def.Motives,EndingRules:def.EndingRules,
  Deductions:def.Deductions,Events:def.Events,CurrentScene:"corridor",
  FoundFacts:map[int]bool{},FoundClues:map[int]bool{},
  SolvedDeductions:map[int]bool{},FoundMotives:map[int]bool{},
  TriggeredEvents:map[int]bool{},OpenDialogues:map[game.DialogueID]bool{},
 }
 move:=func(id string){t.Helper();if err:=game.ChooseSceneOption(s,id);err!=nil{t.Fatalf("%s/%s: %v",s.CurrentScene,id,err)}}
 deduce:=func(id int){t.Helper();for _,d:=range s.Deductions {if d.ID==id {if !game.SolveDeduction(s,d){t.Fatalf("deduction %d unavailable",id)};return}};t.Fatalf("missing deduction %d",id)}
 // First visit an unrelated room; locked evidence cannot be accessed early.
 move("staff")
 if err:=game.ChooseSceneOption(s,"box");err==nil {t.Fatal("weapon revealed before missing-paperweight evidence")}
 move("nikola");move("back")
 move("pantry");move("polish_inventory");move("back")
 move("staff");move("box")
 if err:=game.ChooseSceneOption(s,"cancellation");err==nil {t.Fatal("Sophie cancellation unlocked without note")}
 move("back")
 move("luggage");move("inheritance_copy");move("back")
 move("archive");move("fraud");move("family");move("back")
 move("salon")
 if err:=game.StartConversation(s,"boris_intro");err!=nil {t.Fatal(err)}
 if _,err:=game.ChooseConversationOption(s,"leave");err!=nil {t.Fatal(err)}
 if err:=game.StartConversation(s,"boris_intro");err!=nil {t.Fatal(err)}
 if _,err:=game.ChooseConversationOption(s,"leave");err!=nil {t.Fatal(err)}
 move("back")
 move("victor");move("ribbon");move("watch");move("note");move("back")
 move("staff");move("cancellation");move("back")
 move("vestibule");move("guard_log");move("time_notice");move("back")
 for _,id:=range []int{1,2,4,5} {deduce(id)}
 if !s.FoundFacts[10] || !s.FoundFacts[11] || !s.FoundFacts[15] || !s.FoundMotives[1] {t.Fatal("alternate route lost mandatory evidence")}
 move("salon")
 if err:=game.StartConversation(s,"boris_intro");err!=nil {t.Fatal(err)}
 view,err:=game.CurrentConversationView(s);if err!=nil {t.Fatal(err)}
 available:=false
 for _,choice:=range view.Choices {if choice.ID=="clock_confront" {available=true}}
 if !available {t.Fatal("clock confrontation inaccessible")}
 if _,err:=game.ChooseConversationOption(s,"clock_confront");err!=nil {t.Fatal(err)}
 if _,err:=game.ChooseConversationOption(s,"record");err!=nil {t.Fatal(err)}
 deduce(6)
 for _,id:=range def.Solution.FactIDs {if !s.FoundFacts[id] {t.Fatalf("missing proof %d",id)}}
 if len(s.TriggeredEvents)!=0 {t.Fatal("optional event unexpectedly needed")}
 ending,err:=game.SubmitAccusation(s,game.Accusation{SuspectID:1,MotiveID:1,FactIDs:[]int{10,11,12}})
 if err!=nil {t.Fatal(err)}
 if ending.Type!="win" {t.Fatalf("unexpected verdict %q",ending.Type)}
}

func TestExpress1928LockedChoicesStayHiddenInSceneView(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 s:=&game.GameState{Scenes:def.Scenes,CurrentScene:"staff",FoundFacts:map[int]bool{}}
 view,err:=game.CurrentSceneView(s);if err!=nil {t.Fatal(err)}
 for _,choice:=range view.Choices {
  if choice.ID=="box" || choice.ID=="cancellation" {t.Fatalf("locked choice %q leaked into scene view",choice.ID)}
 }
}

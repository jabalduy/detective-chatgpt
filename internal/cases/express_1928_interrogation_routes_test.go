package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

func TestExpress1928NewEvidenceInterrogations(t *testing.T){
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,issue:=range ValidateCaseIntegrity(def) {t.Error(issue)}
 cases:=[]struct{
  name,scene,conversation,choice,answer string
  evidence []int
  resultFact int
 }{
  {"Marta family","staff","marta_intro","brother","record_brother",[]int{6,7},0},
  {"Boris cancellation","salon","boris_intro","sophie_cancel","record_sophie_cancel",[]int{13,14},15},
  {"Klara fraud","salon","klara_intro","fraud_record","record_fraud_record",[]int{6},0},
  {"Rosa time","salon","rosa_intro","clock_doubt","record_clock_doubt",[]int{5},9},
 }
 for _,tc:=range cases {t.Run(tc.name,func(t *testing.T){
  s:=&game.GameState{Scenes:def.Scenes,Conversations:def.Conversations,Characters:def.Characters,
   CurrentScene:tc.scene,FoundFacts:map[int]bool{},FoundMotives:map[int]bool{}}
  if err:=game.StartConversation(s,tc.conversation);err!=nil {t.Fatal(err)}
  view,err:=game.CurrentConversationView(s);if err!=nil {t.Fatal(err)}
  for _,c:=range view.Choices {if c.ID==tc.choice {t.Fatal("locked interrogation option visible without evidence")}}
  if _,err:=game.ChooseConversationOption(s,"leave");err!=nil {t.Fatal(err)}
  for _,id:=range tc.evidence {s.FoundFacts[id]=true}
  if err:=game.StartConversation(s,tc.conversation);err!=nil {t.Fatal(err)}
  view,err=game.CurrentConversationView(s);if err!=nil {t.Fatal(err)}
  found:=false
  for _,c:=range view.Choices {if c.ID==tc.choice {found=true}}
  if !found {t.Fatalf("option %q still locked after evidence",tc.choice)}
  if _,err:=game.ChooseConversationOption(s,tc.choice);err!=nil {t.Fatal(err)}
  if _,err:=game.ChooseConversationOption(s,tc.answer);err!=nil {t.Fatal(err)}
  if s.ActiveConversation!="" {t.Fatal("conversation not closed")}
  if tc.resultFact!=0 && !s.FoundFacts[tc.resultFact] {t.Fatalf("fact %d not discovered",tc.resultFact)}
 })}
}

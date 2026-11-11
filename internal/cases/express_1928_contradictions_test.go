package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

func TestExpress1928ContradictionInterrogations(t *testing.T){
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,issue:=range ValidateCaseIntegrity(def) {t.Error(issue)}
 cases:=[]struct{
  name,scene,conversation,choice,finish string
  evidence []int
 }{
  {"Marta ribbon and route","staff","marta_intro","ribbon_vs_route","finish_ribbon_vs_route",[]int{1,8}},
  {"Boris watch and Rosa","salon","boris_intro","watch_vs_rosa","finish_watch_vs_rosa",[]int{4,5,9}},
  {"Klara and Pavel","salon","klara_intro","clara_vs_pavel","finish_clara_vs_pavel",[]int{6,7}},
  {"Sophie cancelled meeting","salon","sophie_intro","note_vs_cancel","finish_note_vs_cancel",[]int{13,14}},
 }
 for _,tc:=range cases {t.Run(tc.name,func(t *testing.T){
  s:=&game.GameState{Scenes:def.Scenes,Conversations:def.Conversations,Characters:def.Characters,
   CurrentScene:tc.scene,FoundFacts:map[int]bool{},FoundMotives:map[int]bool{}}
  if err:=game.StartConversation(s,tc.conversation);err!=nil {t.Fatal(err)}
  view,err:=game.CurrentConversationView(s);if err!=nil {t.Fatal(err)}
  for _,c:=range view.Choices {if c.ID==tc.choice {t.Fatal("contradiction visible before evidence")}}
  if _,err:=game.ChooseConversationOption(s,"leave");err!=nil {t.Fatal(err)}
  for _,id:=range tc.evidence {s.FoundFacts[id]=true}
  if err:=game.StartConversation(s,tc.conversation);err!=nil {t.Fatal(err)}
  view,err=game.CurrentConversationView(s);if err!=nil {t.Fatal(err)}
  found:=false
  for _,c:=range view.Choices {if c.ID==tc.choice {found=true}}
  if !found {t.Fatal("contradiction unavailable after evidence")}
  if _,err:=game.ChooseConversationOption(s,tc.choice);err!=nil {t.Fatal(err)}
  if _,err:=game.ChooseConversationOption(s,tc.finish);err!=nil {t.Fatal(err)}
  if s.ActiveConversation!="" {t.Fatal("conversation did not finish")}
 })}
}

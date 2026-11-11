package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

func TestExpress1928ExpandedSceneRoutes(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 if len(def.Scenes)!=8 {t.Fatalf("want eight scenes, got %d",len(def.Scenes))}
 for _,issue:=range ValidateCaseIntegrity(def) {t.Error(issue)}
 for _,tc:=range []struct{
  name string
  scene string
  choice string
  facts []int
  want int
 }{
  {"Nikola via pantry","pantry","tea_log",nil,8},
  {"Boris via vestibule","vestibule","guard_log",nil,4},
  {"Pavel photo via luggage","luggage","family_photo",[]int{6},7},
  {"Rosa time cross-check via vestibule","vestibule","time_notice",[]int{5},9},
 }{
  t.Run(tc.name,func(t *testing.T){
   s:=&game.GameState{Scenes:def.Scenes,CurrentScene:"corridor",FoundFacts:map[int]bool{}}
   if err:=game.ChooseSceneOption(s,tc.scene);err!=nil {t.Fatal(err)}
   for _,id:=range tc.facts {s.FoundFacts[id]=true}
   if err:=game.ChooseSceneOption(s,tc.choice);err!=nil {t.Fatal(err)}
   if !s.FoundFacts[tc.want] {t.Fatalf("fact %d not found",tc.want)}
   if err:=game.ChooseSceneOption(s,"back");err!=nil {t.Fatal(err)}
   if s.CurrentScene!="corridor" {t.Fatalf("cannot return from %s",tc.scene)}
  })
 }
}

func TestExpress1928AlternativeEvidenceRemainsGated(t *testing.T){
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,tc:=range []struct{scene,choice string}{
  {"luggage","family_photo"},
  {"vestibule","time_notice"},
 }{
  t.Run(tc.scene,func(t *testing.T){
   s:=&game.GameState{Scenes:def.Scenes,CurrentScene:tc.scene,FoundFacts:map[int]bool{}}
   if err:=game.ChooseSceneOption(s,tc.choice);err==nil {t.Fatal("evidence prerequisite bypassed")}
  })
 }
}

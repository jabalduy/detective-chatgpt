package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

func TestExpress1928AlternativeSuspectMotives(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,issue:=range ValidateCaseIntegrity(def) {t.Error(issue)}
 motives:=map[int]game.Motive{}
 for _,m:=range def.Motives {motives[m.ID]=m}
 for _,tc:=range []struct{id,suspect int;required []int}{
  {3,5,[]int{6,20}},
  {4,6,[]int{21,22}},
 }{
  m,ok:=motives[tc.id]
  if !ok || m.SuspectID!=tc.suspect {t.Errorf("motive %d not assigned to suspect %d",tc.id,tc.suspect)}
  found:=false
  for _,d:=range def.Deductions {
   for _,e:=range d.Effects {
    if e.Type!="discover_motive" || e.ID!=tc.id {continue}
    found=true
    needed:=map[int]bool{}
    for _,id:=range d.Required {needed[id]=true}
    for _,id:=range tc.required {if !needed[id] {t.Errorf("motive %d missing prerequisite %d",tc.id,id)}}
   }
  }
  if !found {t.Errorf("motive %d has no deduction",tc.id)}
 }
 if def.Solution.SuspectID!=1 || def.Solution.MotiveID!=1 {
  t.Fatal("alternative suspects changed canonical solution")
 }
}

func TestExpress1928FalseMotiveEvidenceIsDiscoverable(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 sources:=map[int]bool{}
 for _,scene:=range def.Scenes {
  for _,choice:=range scene.Choices {
   for _,effect:=range choice.Effects {if effect.Type=="discover_fact" {sources[effect.ID]=true}}
  }
 }
 for _,conversation:=range def.Conversations {
  for _,node:=range conversation.Nodes {
   for _,choice:=range node.Choices {
    for _,effect:=range choice.Effects {if effect.Type=="discover_fact" {sources[effect.ID]=true}}
   }
  }
 }
 for _,id:=range []int{20,21,22} {
  if !sources[id] {t.Errorf("false-suspect evidence %d has no non-event source",id)}
 }
}

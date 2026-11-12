package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

func TestExpress1928GrafAndSophieFairPlay(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,issue:=range ValidateCaseIntegrity(def) {t.Error(issue)}
 if def.Solution.SuspectID!=1 {t.Fatal("canonical killer changed")}
 wanted:=map[int]int{23:4,24:3}
 for factID,characterID:=range wanted {
  found:=false
  for _,conv:=range def.Conversations {
   if conv.CharacterID!=characterID {continue}
   for _,node:=range conv.Nodes {
    for _,choice:=range node.Choices {
     for _,effect:=range choice.Effects {
      if effect.Type=="discover_fact" && effect.ID==factID {found=true}
     }
    }
   }
  }
  if !found {t.Errorf("fact %d not discoverable from character %d",factID,characterID)}
 }
 foundMotive:=false
 for _,deduction:=range def.Deductions {
  for _,effect:=range deduction.Effects {
   if effect.Type!="discover_motive" || effect.ID!=5 {continue}
   foundMotive=true
   facts:=map[int]bool{}
   for _,id:=range deduction.Required {facts[id]=true}
   if !facts[6] || !facts[23] {t.Error("Graf motive lacks documents and confession")}
  }
 }
 if !foundMotive {t.Fatal("Graf alternative motive deduction missing")}
 var graf *game.Conversation
 for i:=range def.Conversations {
  if def.Conversations[i].CharacterID==4 {graf=&def.Conversations[i]}
 }
 if graf==nil {t.Fatal("Graf interrogation missing")}
 gated:=false
 for _,node:=range graf.Nodes {
  if node.ID!=graf.StartNode {continue}
  for _,choice:=range node.Choices {
   if choice.ID!="hidden_letter" {continue}
   for _,req:=range choice.Requirements {
    if req.Type=="fact_found" && req.ID==6 {gated=true}
   }
  }
 }
 if !gated {t.Fatal("Graf letter interrogation is not evidence gated")}
}

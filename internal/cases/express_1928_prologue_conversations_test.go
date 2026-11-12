package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

func TestExpressPrologueConversationsAreSafeAndPlayable(t *testing.T) {
 c,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 want:=map[string]int{"prologue_marta_service":1,"prologue_sophie_salon":3,"prologue_boris_duty":2,"prologue_graf_tea":4,"prologue_klara_salon":5,"prologue_max_salon":6,"prologue_rosa_salon":7,"prologue_leo_archive":8}
 for _,conversation:=range c.Conversations {
  id,ok:=want[conversation.ID]
  if !ok {continue}
  if conversation.AvailableAct!=game.ActPrologue {t.Errorf("%s act=%q",conversation.ID,conversation.AvailableAct)}
  if conversation.CharacterID!=id {t.Errorf("%s character=%d",conversation.ID,conversation.CharacterID)}
  nodes:=map[string]bool{}
  for _,node:=range conversation.Nodes {nodes[node.ID]=true}
  if !nodes[conversation.StartNode] {t.Errorf("%s missing start node",conversation.ID)}
  for _,node:=range conversation.Nodes {
   for _,choice:=range node.Choices {
    if len(choice.Effects)!=0 {t.Errorf("%s choice %s changes investigation state",conversation.ID,choice.ID)}
    if !choice.End && !nodes[choice.NextNode] {t.Errorf("%s choice %s has missing next node",conversation.ID,choice.ID)}
   }
  }
  delete(want,conversation.ID)
 }
 for id:=range want {t.Errorf("missing prologue conversation %s",id)}
}

package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

func TestExpressEvidenceConfrontations(t *testing.T) {
 def, err := LoadCase(filepath.Join("..", "..", "cases", "case_express_1928"))
 if err != nil { t.Fatal(err) }
 investigationCount := 0
 for _, c := range def.Conversations {
  if c.AvailableAct == game.ActPrologue { continue }
  investigationCount++
  if c.StartNode == "" { t.Errorf("%s missing start", c.ID);continue }
  nodes := make(map[string]bool)
  for _, n := range c.Nodes {nodes[n.ID]=true}
  if !nodes[c.StartNode] {t.Errorf("%s missing initial node",c.ID)}
  guarded := false
  for _, n := range c.Nodes {
   for _, ch := range n.Choices {
    if !ch.End && !nodes[ch.NextNode] {t.Errorf("%s/%s points to missing node %q",c.ID,n.ID,ch.NextNode)}
    if ch.ID=="evidence" {
     if len(ch.Requirements)!=1 || ch.Requirements[0].Type!="fact_found" {t.Errorf("%s confrontation not evidence gated",c.ID)}
     guarded=true
    }
   }
  }
  if !guarded {t.Errorf("%s has no evidence confrontation",c.ID)}
 }
 if investigationCount != 8 {t.Fatalf("want 8 investigation conversations, got %d",investigationCount)}
}

func TestExpressConfrontationRequiresEvidence(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,c:=range def.Conversations {
  for _,n:=range c.Nodes {
   for _,ch:=range n.Choices {
    if ch.ID!="evidence" {continue}
    // This test deliberately checks schema, not GameState runtime setup.
    if len(ch.Requirements)!=1 || ch.Requirements[0].ID<1 {
     t.Errorf("%s has invalid evidence prerequisite",c.ID)
    }
   }
  }
 }
 _ = game.Requirement{}
}

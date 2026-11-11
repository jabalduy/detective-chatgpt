package cases

import (
 "path/filepath"
 "testing"
)

func TestExpressClockEvidenceChain(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 facts:=map[int]bool{}
 for _,f:=range def.Facts {facts[f.ID]=true}
 for _,id:=range []int{4,5,9,12,13,14,15,16} {
  if !facts[id] {t.Errorf("missing fact %d",id)}
 }
 deductions:=map[int][]int{}
 for _,d:=range def.Deductions {deductions[d.ID]=d.Required}
 if !sameIntSet(deductions[5],[]int{13,14}) {t.Error("Sophie contradiction must require note and cancellation")}
 if !sameIntSet(deductions[6],[]int{5,9,16}) {t.Error("confirmed clock tampering must require Boris admission")}
 for _,d:=range def.Deductions {
  if d.ID==3 {
   for _,effect:=range d.Effects {
    if effect.Type=="discover_fact" && effect.ID==12 {t.Error("circumstantial timing must not prove tampering")}
   }
  }
 }
 found:=false
 for _,c:=range def.Conversations {
  if c.ID!="boris_intro" {continue}
  for _,node:=range c.Nodes {
   for _,ch:=range node.Choices {
    if ch.ID!="clock_confront" {continue}
    required:=map[int]bool{}
    for _,r:=range ch.Requirements {if r.Type=="fact_found" {required[r.ID]=true}}
    if !required[4] || !required[9] || !required[15] {t.Error("Boris confrontation missing prerequisites")}
    found=true
   }
  }
 }
 if !found {t.Fatal("missing Boris confrontation")}
}
func sameIntSet(a,b []int)bool{
 if len(a)!=len(b){return false}
 seen:=map[int]bool{}
 for _,n:=range a{seen[n]=true}
 for _,n:=range b{if !seen[n]{return false}}
 return true
}

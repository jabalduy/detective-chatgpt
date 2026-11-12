package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

// This measures authored action costs, not real player deliberation time.
// It provides a deterministic lower bound for a successful investigation.
func TestExpress1928CriticalPathTimeBudget(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 s:=&game.GameState{
  Scenes:def.Scenes,Conversations:def.Conversations,Characters:def.Characters,
  Deductions:def.Deductions,CurrentScene:"corridor",
  FoundFacts:map[int]bool{},FoundMotives:map[int]bool{},
  SolvedDeductions:map[int]bool{},
 }
 elapsed:=0
 step:=func(id string){
  t.Helper()
  scene:=game.SceneByID(s,s.CurrentScene)
  if scene==nil {t.Fatalf("missing scene %q",s.CurrentScene)}
  found:=false
  for _,c:=range scene.Choices {if c.ID==id {elapsed+=c.TimeCost;found=true;break}}
  if !found {t.Fatalf("missing choice %s in %s",id,s.CurrentScene)}
  if err:=game.ChooseSceneOption(s,id);err!=nil {t.Fatal(err)}
 }
 for _,id:=range []string{"victor","ribbon","paperweight","watch","note","back","staff","box","nikola","cancellation","back","salon","boris","rosa","back","archive","fraud","family"} {step(id)}
 for _,id:=range []int{1,2,4,5} {
  found:=false
  for _,d:=range s.Deductions {if d.ID==id {
   if !game.SolveDeduction(s,d) {t.Fatalf("cannot solve deduction %d",id)}
   elapsed+=d.TimeCost
   found=true;break
  }}
  if !found {t.Fatalf("missing deduction %d",id)}
 }
 for _,id:=range []string{"back","salon"} {step(id)}
 if err:=game.StartConversation(s,"boris_intro");err!=nil {t.Fatal(err)}
 for _,id:=range []string{"clock_confront","record"} {
  conversation:=s.ActiveConversation
  nodeID:=s.ConversationNode
  found:=false
  for _,c:=range s.Conversations {if c.ID==conversation {
   for _,node:=range c.Nodes {if node.ID==nodeID {
    for _,choice:=range node.Choices {if choice.ID==id {elapsed+=choice.TimeCost;found=true}}
   }}
  }}
  if !found {t.Fatalf("missing conversation choice %s",id)}
  if _,err:=game.ChooseConversationOption(s,id);err!=nil {t.Fatal(err)}
 }
 for _,d:=range s.Deductions {if d.ID==6 {
  if !game.SolveDeduction(s,d) {t.Fatal("final deduction unavailable")}
  elapsed+=d.TimeCost
 }}
 for _,id:=range def.Solution.FactIDs {if !s.FoundFacts[id] {t.Fatalf("solution fact %d missing",id)}}
 const investigationSeconds=60*60
 if elapsed>=investigationSeconds {t.Fatalf("critical path costs %d seconds, exceeds %d",elapsed,investigationSeconds)}
 t.Logf("critical path authored action costs: %d seconds (%.1f minutes); remaining %.1f minutes for reading, discussion and optional content",elapsed,float64(elapsed)/60,float64(investigationSeconds-elapsed)/60)
}

func TestExpress1928TimedEventsHaveNonBlockingDefaults(t *testing.T){
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 if len(def.Events)!=5 {t.Fatalf("expected 5 timed events, got %d",len(def.Events))}
 for _,event:=range def.Events {
  if event.TimeLimit!=10 {t.Errorf("event %d: want 10s decision window, got %d",event.ID,event.TimeLimit)}
  found:=false
  for _,option:=range event.Options {if option.ID==event.DefaultOptionID {
   found=true
   if len(option.Requirements)>0 {t.Errorf("event %d: default has requirements",event.ID)}
  }}
  if !found {t.Errorf("event %d: no default option",event.ID)}
 }
}

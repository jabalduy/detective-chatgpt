package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

// Fair-play invariant: optional timed scenes must not be the sole source of
// any evidence needed for the canonical solution.
func TestExpress1928OptionalEventsDoNotGateSolution(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 required:=map[int]bool{}
 for _,id:=range def.Solution.FactIDs {required[id]=true}
 sceneFacts:=map[int]bool{}
 for _,scene:=range def.Scenes {
  for _,choice:=range scene.Choices {
   for _,effect:=range choice.Effects {
    if effect.Type=="discover_fact" {sceneFacts[effect.ID]=true}
   }
  }
 }
 eventOnly:=map[int]bool{}
 for _,event:=range def.Events {
  if event.TimeLimit!=10 {t.Errorf("event %d has %d-second limit, expected 10",event.ID,event.TimeLimit)}
  for _,option:=range event.Options {
   for _,effect:=range option.Effects {
    if effect.Type=="discover_fact" && !sceneFacts[effect.ID] {eventOnly[effect.ID]=true}
   }
  }
 }
 for _,deduction:=range def.Deductions {
  for _,id:=range deduction.Required {
   if eventOnly[id] {t.Errorf("deduction %d requires event-only fact %d",deduction.ID,id)}
  }
 }
 for id:=range required {
  if eventOnly[id] {t.Errorf("mandatory proof %d is locked behind timed event",id)}
 }
}

func TestExpress1928EverySuspectCanBeQuestionedAndLeave(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 conversations:=map[int]game.Conversation{}
 for _,conversation:=range def.Conversations {
  if _,exists:=conversations[conversation.CharacterID];exists {
   t.Errorf("duplicate suspect conversation for character %d",conversation.CharacterID)
  }
  conversations[conversation.CharacterID]=conversation
 }
 for _,suspect:=range def.Suspects {
  conversation,ok:=conversations[suspect.SusID]
  if !ok {t.Errorf("suspect %d has no conversation",suspect.SusID);continue}
  hasExit:=false
  hasQuestion:=false
  for _,node:=range conversation.Nodes {
   if node.ID!=conversation.StartNode {continue}
   for _,choice:=range node.Choices {
    if choice.End && len(choice.Requirements)==0 {hasExit=true}
    if !choice.End && len(choice.Requirements)==0 {hasQuestion=true}
   }
  }
  if !hasExit {t.Errorf("suspect %d cannot exit initial interrogation without prerequisites",suspect.SusID)}
  if !hasQuestion {t.Errorf("suspect %d has no available opening question",suspect.SusID)}
 }
}

func TestExpress1928ClockConfessionNeedsCorroboration(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 var boris *game.Conversation
 for i:=range def.Conversations {
  if def.Conversations[i].CharacterID==2 {boris=&def.Conversations[i];break}
 }
 if boris==nil {t.Fatal("Boris conversation missing")}
 found:=false
 for _,node:=range boris.Nodes {
  if node.ID!=boris.StartNode {continue}
  for _,choice:=range node.Choices {
   if choice.ID!="clock_confront" {continue}
   found=true
   prereqs:=map[int]bool{}
   for _,req:=range choice.Requirements {
    if req.Type=="fact_found" {prereqs[req.ID]=true}
   }
   for _,id:=range []int{4,9,15} {
    if !prereqs[id] {t.Errorf("clock confession can be requested without prerequisite fact %d",id)}
   }
  }
 }
 if !found {t.Fatal("Boris has no clock confrontation")}
}

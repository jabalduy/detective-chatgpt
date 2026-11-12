package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

// This audit follows the authored clue chain, not merely the JSON schema.
// A final proof must require a deduction and every input to that deduction
// must have a source outside optional ten-second events.
func TestExpress1928FinalProofHasIndependentSources(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,issue:=range ValidateCaseIntegrity(def) {t.Error(issue)}
 direct:=map[int]bool{}
 for _,scene:=range def.Scenes {
  for _,choice:=range scene.Choices {
   for _,effect:=range choice.Effects {
    if effect.Type=="discover_fact" {direct[effect.ID]=true}
   }
  }
 }
 for _,conversation:=range def.Conversations {
  for _,node:=range conversation.Nodes {
   for _,choice:=range node.Choices {
    for _,effect:=range choice.Effects {
     if effect.Type=="discover_fact" {direct[effect.ID]=true}
    }
   }
  }
 }
 produced:=map[int]game.Deduction{}
 for _,deduction:=range def.Deductions {
  for _,effect:=range deduction.Effects {
   if effect.Type=="discover_fact" {produced[effect.ID]=deduction}
  }
 }
 var visit func(int,map[int]bool)
 visit=func(id int,stack map[int]bool){
  if direct[id] {return}
  if stack[id] {t.Errorf("circular deduction dependency at fact %d",id);return}
  deduction,ok:=produced[id]
  if !ok {t.Errorf("fact %d has no source outside timed events",id);return}
  stack[id]=true
  for _,required:=range deduction.Required {visit(required,stack)}
  delete(stack,id)
 }
 for _,id:=range def.Solution.FactIDs {
  if direct[id] {t.Errorf("final proof %d bypasses deduction",id)}
  visit(id,map[int]bool{})
 }
 // The clock chain includes a statement by Boris. It must not be
 // awarded merely by walking into a scene.
 if direct[16] {
  for _,scene:=range def.Scenes {
   for _,choice:=range scene.Choices {
    for _,effect:=range choice.Effects {
     if effect.Type=="discover_fact" && effect.ID==16 {
      t.Errorf("scene %s gives Boris confession without interrogation",scene.ID)
     }
    }
   }
  }
 }
}

func TestExpress1928TimedEventsCannotAwardVerdictOrMotive(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 mandatory:=map[int]bool{}
 for _,id:=range def.Solution.FactIDs {mandatory[id]=true}
 for _,event:=range def.Events {
  for _,option:=range event.Options {
   for _,effect:=range option.Effects {
    if effect.Type=="discover_motive" {t.Errorf("event %d grants motive",event.ID)}
    if effect.Type=="discover_fact" && mandatory[effect.ID] {
     t.Errorf("event %d grants mandatory proof %d",event.ID,effect.ID)
    }
   }
  }
 }
}

func TestExpress1928AllEndingRulesHaveValidTarget(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 suspects:=map[int]bool{}
 for _,s:=range def.Suspects {suspects[s.SusID]=true}
 for _,rule:=range def.EndingRules {
  if rule.ID!="win" && rule.ID!="fail" && rule.ID!="unsolved" {
   t.Errorf("unknown ending outcome %q",rule.ID)
  }
  if rule.SuspectID!=0 && !suspects[rule.SuspectID] {
   t.Errorf("ending rule %q targets unknown suspect %d",rule.Title,rule.SuspectID)
  }
  if rule.Title=="" || rule.Text=="" {
   t.Errorf("ending rule %q has empty narrative",rule.ID)
  }
 }
}

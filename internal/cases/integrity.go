package cases

import (
 "fmt"
 "sort"

 "detective/internal/game"
)

// ValidateCaseIntegrity checks cross-file references that JSON unmarshalling
// and per-subsystem narrative validation cannot detect.
func ValidateCaseIntegrity(def *CaseDefinition) []error {
 if def==nil {return []error{fmt.Errorf("nil case definition")}}
 var issues []error
 issues=append(issues,game.ValidateNarrative(def.Scenes,def.Conversations)...)
 issues=append(issues,game.ValidateEvents(def.Events)...)
 facts:=map[int]bool{}
 motives:=map[int]bool{}
 suspects:=map[int]bool{}
 scenes:=map[string]bool{}
 for _,f:=range def.Facts {if facts[f.ID] {issues=append(issues,fmt.Errorf("duplicate fact %d",f.ID))};facts[f.ID]=true}
 for _,m:=range def.Motives {motives[m.ID]=true}
 for _,s:=range def.Suspects {suspects[s.SusID]=true}
 for _,s:=range def.Scenes {scenes[s.ID]=true}
 checkRequirements:=func(where string,reqs []game.Requirement){
  for _,r:=range reqs {
   if r.Type=="fact_found" && !facts[r.ID] {issues=append(issues,fmt.Errorf("%s requires unknown fact %d",where,r.ID))}
   if (r.Type=="character_present" || r.Type=="character_absent" || r.Type=="character_trust_at_least" || r.Type=="character_trust_below") && !suspects[r.ID] {
    issues=append(issues,fmt.Errorf("%s references unknown character %d",where,r.ID))
   }
  }
 }
 checkEffects:=func(where string,effects []game.Effect){
  for _,e:=range effects {
   if e.Type=="discover_fact" && !facts[e.ID] {issues=append(issues,fmt.Errorf("%s discovers unknown fact %d",where,e.ID))}
   if e.Type=="discover_motive" && !motives[e.ID] {issues=append(issues,fmt.Errorf("%s discovers unknown motive %d",where,e.ID))}
   if (e.Type=="character_trust_up" || e.Type=="character_trust_down" || e.Type=="character_leave" || e.Type=="character_appear") && !suspects[e.ID] {
    issues=append(issues,fmt.Errorf("%s affects unknown character %d",where,e.ID))
   }
  }
 }
 for _,s:=range def.Scenes {
  for _,c:=range s.Choices {
   where:=fmt.Sprintf("scene %s choice %s",s.ID,c.ID)
   checkRequirements(where,c.Requirements);checkEffects(where,c.Effects)
  }
 }
 for _,c:=range def.Conversations {
  if !suspects[c.CharacterID] {issues=append(issues,fmt.Errorf("conversation %s references unknown character %d",c.ID,c.CharacterID))}
  if character,ok:=def.Characters[c.CharacterID];!ok || !scenes[character.Location] {
   issues=append(issues,fmt.Errorf("conversation %s character cannot be placed in a scene",c.ID))
  }
  for _,node:=range c.Nodes {
   for _,choice:=range node.Choices {
    where:=fmt.Sprintf("conversation %s node %s choice %s",c.ID,node.ID,choice.ID)
    checkRequirements(where,choice.Requirements);checkEffects(where,choice.Effects)
   }
  }
 }
 for _,e:=range def.Events {
  for _,o:=range e.Options {
   where:=fmt.Sprintf("event %d option %d",e.ID,o.ID)
   checkRequirements(where,o.Requirements);checkEffects(where,o.Effects)
  }
 }
 for _,d:=range def.Deductions {
  for _,id:=range d.Required {if !facts[id] {issues=append(issues,fmt.Errorf("deduction %d requires unknown fact %d",d.ID,id))}}
  checkRequirements(fmt.Sprintf("deduction %d",d.ID),d.Requirements)
  checkEffects(fmt.Sprintf("deduction %d",d.ID),d.Effects)
 }
 for _,r:=range def.EndingRules {checkRequirements("ending rule "+r.ID,r.Requirements)}
 if !suspects[def.Solution.SuspectID] {issues=append(issues,fmt.Errorf("solution has unknown suspect %d",def.Solution.SuspectID))}
 if !motives[def.Solution.MotiveID] {issues=append(issues,fmt.Errorf("solution has unknown motive %d",def.Solution.MotiveID))}
 for _,id:=range def.Solution.FactIDs {if !facts[id] {issues=append(issues,fmt.Errorf("solution requires unknown fact %d",id))}}
 sort.Slice(issues,func(i,j int)bool{return issues[i].Error()<issues[j].Error()})
 return issues
}

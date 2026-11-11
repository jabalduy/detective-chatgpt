package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

// A static reachability audit complements the playable route tests: every
// scene must be connected, and all mandatory proof must have a non-event
// acquisition chain. It deliberately ignores optional timed events.
func TestExpress1928NarrativeReachabilityAudit(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,issue:=range ValidateCaseIntegrity(def) {t.Error(issue)}
 scenes:=map[string]game.Scene{}
 for _,scene:=range def.Scenes {scenes[scene.ID]=scene}
 reached:=map[string]bool{"corridor":true}
 for changed:=true;changed; {
  changed=false
  for id:=range reached {
   for _,choice:=range scenes[id].Choices {
    if choice.NextScene!="" && !reached[choice.NextScene] {
     reached[choice.NextScene]=true;changed=true
    }
   }
  }
 }
 for _,scene:=range def.Scenes {
  if !reached[scene.ID] {t.Errorf("scene %q is unreachable from corridor",scene.ID)}
 }
 // Fixed-point closure: treat all reachable scene choices as potential
 // sources, but only when their fact prerequisites have already been found.
 // This is an over-approximation, not a substitute for gameplay tests.
 facts:=map[int]bool{}
 motives:=map[int]bool{}
 apply:=func(effects []game.Effect)bool{
  changed:=false
  for _,e:=range effects {
   if e.Type=="discover_fact" && !facts[e.ID] {facts[e.ID]=true;changed=true}
   if e.Type=="discover_motive" && !motives[e.ID] {motives[e.ID]=true;changed=true}
  }
  return changed
 }
 available:=func(reqs []game.Requirement)bool{
  for _,r:=range reqs {
   if r.Type=="fact_found" && !facts[r.ID] {return false}
  }
  return true
 }
 for changed:=true;changed; {
  changed=false
  for _,scene:=range def.Scenes {
   if !reached[scene.ID] {continue}
   for _,choice:=range scene.Choices {
    if available(choice.Requirements) && apply(choice.Effects) {changed=true}
   }
  }
  for _,conversation:=range def.Conversations {
   for _,node:=range conversation.Nodes {
    for _,choice:=range node.Choices {
     if available(choice.Requirements) && apply(choice.Effects) {changed=true}
    }
   }
  }
  for _,d:=range def.Deductions {
   ok:=available(d.Requirements)
   for _,id:=range d.Required {if !facts[id] {ok=false}}
   if ok && apply(d.Effects) {changed=true}
  }
 }
 for _,id:=range def.Solution.FactIDs {
  if !facts[id] {t.Errorf("mandatory solution fact %d has no reachable acquisition chain",id)}
 }
 if !motives[def.Solution.MotiveID] {t.Errorf("solution motive %d cannot be discovered",def.Solution.MotiveID)}
 for _,id:=range []int{13,14,15,16} {
  if !facts[id] {t.Errorf("clock confession prerequisite fact %d is unreachable",id)}
 }
}

func TestExpress1928EverySceneHasExit(t *testing.T){
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,scene:=range def.Scenes {
  if scene.ID=="corridor" {continue}
  hasExit:=false
  for _,choice:=range scene.Choices {
   if choice.NextScene=="corridor" && len(choice.Requirements)==0 {hasExit=true}
  }
  if !hasExit {t.Errorf("scene %q has no unconditional return to corridor",scene.ID)}
 }
}

package game

import (
 "fmt"
 "sort"
)

// SceneReachabilityReport checks whether evidence can be obtained by moving
// through scenes and taking scene actions. It intentionally ignores dialogue,
// timed events and other producers; those require a full game simulation.
type SceneReachabilityReport struct {
 ReachableScenes []string
 MissingSolutionFacts []int
}

// AnalyzeSceneReachability repeatedly discovers accessible scenes and positive
// evidence effects. This is an optimistic analysis: choices may be revisited,
// and no irreversible state or clock costs are simulated.
func AnalyzeSceneReachability(s *GameState) SceneReachabilityReport {
 reached:=map[string]bool{}
 if SceneByID(s,s.CurrentScene)!=nil {reached[s.CurrentScene]=true}
 clues:=map[int]bool{}
 facts:=map[int]bool{}
 for id,v:=range s.FoundClues {if v{clues[id]=true}}
 for id,v:=range s.FoundFacts {if v{facts[id]=true}}
 canUse:=func(reqs []Requirement)bool{
  for _,r:=range reqs {
   switch r.Type {
   case "clue_found":if !clues[r.ID]{return false}
   case "fact_found":if !facts[r.ID]{return false}
   default:return false // Unknown prerequisites cannot be assumed satisfied.
   }
  }
  return true
 }
 for changed:=true;changed;{
  changed=false
  for _,scene:=range s.Scenes {
   if !reached[scene.ID]{continue}
   for _,choice:=range scene.Choices {
    if !canUse(choice.Requirements){continue}
    if choice.NextScene!="" && SceneByID(s,choice.NextScene)!=nil && !reached[choice.NextScene]{
     reached[choice.NextScene]=true;changed=true
    }
    for _,e:=range choice.Effects {
     switch e.Type {
     case "discover_clue":if !clues[e.ID]{clues[e.ID]=true;changed=true}
     case "discover_fact":if !facts[e.ID]{facts[e.ID]=true;changed=true}
     }
    }
   }
  }
 }
 result:=SceneReachabilityReport{ReachableScenes:[]string{},MissingSolutionFacts:[]int{}}
 for id:=range reached {result.ReachableScenes=append(result.ReachableScenes,id)}
 sort.Strings(result.ReachableScenes)
 for _,id:=range s.Case.Solution.FactIDs {
  if !facts[id]{result.MissingSolutionFacts=append(result.MissingSolutionFacts,id)}
 }
 sort.Ints(result.MissingSolutionFacts)
 return result
}

// SceneReachabilityDiagnostics returns actionable warnings. A missing fact
// might still be obtainable through conversations, examinations or events.
func SceneReachabilityDiagnostics(s *GameState)[]string{
 r:=AnalyzeSceneReachability(s)
 warnings:=[]string{}
 for _,scene:=range s.Scenes {
  found:=false
  for _,id:=range r.ReachableScenes {if id==scene.ID{found=true;break}}
  if !found{warnings=append(warnings,fmt.Sprintf("scene %q unreachable from %q through scene choices",scene.ID,s.CurrentScene))}
 }
 return warnings
}

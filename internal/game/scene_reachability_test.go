package game

import (
 "slices"
 "testing"
)

func TestSceneReachabilityUnlocksLaterRoute(t *testing.T){
 s:=&GameState{
  CurrentScene:"hall",
  Case:CaseInfo{Solution:Solution{FactIDs:[]int{8}}},
  Scenes:[]Scene{
   {ID:"hall",Choices:[]Choice{
    {ID:"door",NextScene:"lab",Action:Action{Requirements:[]Requirement{{Type:"clue_found",ID:4}}}},
    {ID:"search",Action:Action{Effects:[]Effect{{Type:"discover_clue",ID:4}}}},
   }},
   {ID:"lab",Choices:[]Choice{{ID:"inspect",Action:Action{Effects:[]Effect{{Type:"discover_fact",ID:8}}}}}},
  },
 }
 r:=AnalyzeSceneReachability(s)
 if !slices.Contains(r.ReachableScenes,"lab") || len(r.MissingSolutionFacts)!=0{
  t.Fatalf("expected lab and fact to become reachable: %+v",r)
 }
 if s.FoundClues[4] || s.FoundFacts[8] || s.CurrentScene!="hall"{t.Fatal("analysis mutated live game state")}
}
func TestSceneReachabilityDetectsLockedCycle(t *testing.T){
 s:=&GameState{
  CurrentScene:"hall",Case:CaseInfo{Solution:Solution{FactIDs:[]int{8}}},
  Scenes:[]Scene{
   {ID:"hall",Choices:[]Choice{{ID:"door",NextScene:"lab",Action:Action{Requirements:[]Requirement{{Type:"fact_found",ID:8}}}}}},
   {ID:"lab",Choices:[]Choice{{ID:"inspect",Action:Action{Effects:[]Effect{{Type:"discover_fact",ID:8}}}}}},
  },
 }
 r:=AnalyzeSceneReachability(s)
 if slices.Contains(r.ReachableScenes,"lab") || !slices.Contains(r.MissingSolutionFacts,8){
  t.Fatalf("locked cycle missed: %+v",r)
 }
 if len(SceneReachabilityDiagnostics(s))!=1{t.Fatal("expected unreachable scene warning")}
}
func TestSceneReachabilityDoesNotAssumeUnknownRequirements(t *testing.T){
 s:=&GameState{CurrentScene:"hall",Scenes:[]Scene{
  {ID:"hall",Choices:[]Choice{{ID:"door",NextScene:"lab",Action:Action{Requirements:[]Requirement{{Type:"character_present",ID:1}}}}}},
  {ID:"lab"},
 }}
 if slices.Contains(AnalyzeSceneReachability(s).ReachableScenes,"lab"){t.Fatal("unsupported requirement bypassed")}
}

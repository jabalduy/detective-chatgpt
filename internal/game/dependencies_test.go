package game

import (
 "slices"
 "testing"
)

func TestAnalyzeDependenciesReachableChain(t *testing.T){
 s:=&GameState{
  Scenes:[]Scene{{ID:"hall",Choices:[]Choice{{ID:"find",Action:Action{Effects:[]Effect{{Type:"discover_clue",ID:1}}}}}}},
  Examinations:[]Examination{{ID:"lab",Requirements:[]Requirement{{Type:"clue_found",ID:1}},Effects:[]Effect{{Type:"discover_fact",ID:2}}}},
  Case:CaseInfo{Solution:Solution{FactIDs:[]int{2}}},
 }
 result:=AnalyzeDependencies(s)
 if len(result.Unreachable)!=0{t.Fatalf("unexpected unreachable evidence: %v",result.Unreachable)}
}
func TestAnalyzeDependenciesCircularEvidence(t *testing.T){
 s:=&GameState{Scenes:[]Scene{{ID:"hall",Choices:[]Choice{
  {ID:"first",Action:Action{Requirements:[]Requirement{{Type:"fact_found",ID:2}},Effects:[]Effect{{Type:"discover_clue",ID:1}}}},
  {ID:"second",Action:Action{Requirements:[]Requirement{{Type:"clue_found",ID:1}},Effects:[]Effect{{Type:"discover_fact",ID:2}}}},
 }}}}
 result:=AnalyzeDependencies(s)
 if !slices.Contains(result.Unreachable,"clue:1") || !slices.Contains(result.Unreachable,"fact:2"){
  t.Fatalf("circular dependency not detected: %v",result.Unreachable)
 }
}
func TestAnalyzeDependenciesAlternativePath(t *testing.T){
 s:=&GameState{Scenes:[]Scene{{ID:"hall",Choices:[]Choice{
  {ID:"blocked",Action:Action{Requirements:[]Requirement{{Type:"fact_found",ID:2}},Effects:[]Effect{{Type:"discover_clue",ID:1}}}},
  {ID:"alternative",Action:Action{Effects:[]Effect{{Type:"discover_clue",ID:1}}}},
  {ID:"fact",Action:Action{Requirements:[]Requirement{{Type:"clue_found",ID:1}},Effects:[]Effect{{Type:"discover_fact",ID:2}}}},
 }}}}
 if result:=AnalyzeDependencies(s);len(result.Unreachable)!=0{t.Fatalf("alternative path missed: %v",result.Unreachable)}
}
func TestAnalyzeDependenciesMissingProducer(t *testing.T){
 s:=&GameState{Case:CaseInfo{Solution:Solution{FactIDs:[]int{99}}}}
 if result:=AnalyzeDependencies(s);!slices.Contains(result.Unreachable,"fact:99"){
  t.Fatalf("missing solution fact not detected: %v",result.Unreachable)
 }
}

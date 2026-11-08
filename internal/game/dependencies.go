package game

import (
 "fmt"
 "sort"
)

// DependencyReport is an optimistic, order-independent reachability analysis.
// It detects missing producers and circular positive dependencies; it is not a
// substitute for simulating scene movement, timing or mutually exclusive choices.
type DependencyReport struct {
 Unreachable []string
}

type dependencyAction struct {
 requires []string
 produces []string
}

func dependencyRequirement(r Requirement)(string,bool){
 switch r.Type {
 case "clue_found":return fmt.Sprintf("clue:%d",r.ID),true
 case "fact_found":return fmt.Sprintf("fact:%d",r.ID),true
 case "examination_completed":return "examination:"+r.Key,true
 default:return "",false
 }
}
func dependencyEffect(e Effect)(string,bool){
 switch e.Type {
 case "discover_clue":return fmt.Sprintf("clue:%d",e.ID),true
 case "discover_fact":return fmt.Sprintf("fact:%d",e.ID),true
 default:return "",false
 }
}
func dependencyActionFrom(a Action)dependencyAction{
 result:=dependencyAction{}
 for _,r:=range a.Requirements {if key,ok:=dependencyRequirement(r);ok{result.requires=append(result.requires,key)}}
 for _,e:=range a.Effects {if key,ok:=dependencyEffect(e);ok{result.produces=append(result.produces,key)}}
 return result
}

// AnalyzeDependencies finds evidence that cannot be obtained even when all
// possible actions are considered available and repeatable. An unreachable
// item is therefore a strong warning of missing/circular prerequisites.
func AnalyzeDependencies(s *GameState)DependencyReport{
 actions:=make([]dependencyAction,0)
 for _,scene:=range s.Scenes {for _,choice:=range scene.Choices {
  actions=append(actions,dependencyActionFrom(choice.Action))
 }}
 for _,conversation:=range s.Conversations {for _,node:=range conversation.Nodes {for _,choice:=range node.Choices {
  actions=append(actions,dependencyActionFrom(choice.Action))
 }}}
 for _,object:=range s.Objects {actions=append(actions,dependencyActionFrom(object.InspectAction))}
 for _,dialogue:=range s.Dialogues {
  actions=append(actions,dependencyActionFrom(dialogue.Action))
  if dialogue.Confrontation!=nil {actions=append(actions,dependencyActionFrom(dialogue.Confrontation.Action))}
 }
 for _,event:=range s.Events {for _,option:=range event.Options {
  actions=append(actions,dependencyActionFrom(option.Action))
 }}
 for _,p:=range s.EvidencePresentations {
  a:=dependencyActionFrom(Action{Requirements:p.Requirements,Effects:p.Effects})
  switch p.EvidenceType {
  case "clue":a.requires=append(a.requires,fmt.Sprintf("clue:%d",p.EvidenceID))
  case "fact":a.requires=append(a.requires,fmt.Sprintf("fact:%d",p.EvidenceID))
  }
  actions=append(actions,a)
 }
 for _,examination:=range s.Examinations {
  a:=dependencyActionFrom(Action{Requirements:examination.Requirements,Effects:examination.Effects})
  a.produces=append(a.produces,"examination:"+examination.ID)
  actions=append(actions,a)
 }
 available:=make(map[string]bool)
 for id,found:=range s.FoundClues {if found{available[fmt.Sprintf("clue:%d",id)]=true}}
 for id,found:=range s.FoundFacts {if found{available[fmt.Sprintf("fact:%d",id)]=true}}
 for id,done:=range s.ExaminationCompleted {if done{available["examination:"+id]=true}}
 for {
  changed:=false
  for _,a:=range actions {
   ready:=true
   for _,req:=range a.requires {if !available[req]{ready=false;break}}
   if !ready{continue}
   for _,effect:=range a.produces {if !available[effect]{available[effect]=true;changed=true}}
  }
  if !changed{break}
 }
 all:=make(map[string]bool)
 for _,a:=range actions {
  for _,req:=range a.requires {all[req]=true}
  for _,effect:=range a.produces {all[effect]=true}
 }
 for _,id:=range s.Case.Solution.FactIDs {all[fmt.Sprintf("fact:%d",id)]=true}
 result:=DependencyReport{Unreachable:[]string{}}
 for key:=range all {if !available[key]{result.Unreachable=append(result.Unreachable,key)}}
 sort.Strings(result.Unreachable)
 return result
}

package game

import (
 "fmt"
 "slices"
)

type Accusation struct {
 SuspectID int `json:"suspect_id"`
 MotiveID int `json:"motive_id"`
 FactIDs []int `json:"fact_ids"`
}
type GameEnding struct {
 Type string `json:"type"`
 Title string `json:"title"`
 Text string `json:"text"`
 Accusation Accusation `json:"accusation"`
}
type EndingRule struct {
 ID string `json:"id"`
 Title string `json:"title"`
 Text string `json:"text"`
 Requirements []Requirement `json:"requirements,omitempty"`
 Priority int `json:"priority,omitempty"`
}
func SubmitAccusation(s *GameState,a Accusation)(GameEnding,error){
 if s.Ending!=nil{return GameEnding{},fmt.Errorf("game already finished")}
 if s.ActiveEvent!=nil || s.ActiveConversation!="" || s.ActiveDialogue!=nil {return GameEnding{},fmt.Errorf("finish active interaction first")}
 if !s.FoundMotives[a.MotiveID]{return GameEnding{},fmt.Errorf("motive not discovered")}
 motiveValid:=false
 for _,m:=range s.Motives {if m.ID==a.MotiveID && m.SuspectID==a.SuspectID {motiveValid=true;break}}
 if !motiveValid{return GameEnding{},fmt.Errorf("motive does not match suspect")}
 seenFacts:=make(map[int]bool,len(a.FactIDs))
 for _,id:=range a.FactIDs {
  if seenFacts[id] {return GameEnding{},fmt.Errorf("duplicate fact %d",id)}
  seenFacts[id]=true
  if !s.FoundFacts[id]{return GameEnding{},fmt.Errorf("fact %d not discovered",id)}
 }
 solution:=s.Case.Solution
 outcome:="win"
 if a.SuspectID!=solution.SuspectID {outcome="fail"
 } else if a.MotiveID!=solution.MotiveID {outcome="unsolved"
 } else {for _,id:=range solution.FactIDs {if !slices.Contains(a.FactIDs,id){outcome="unsolved";break}}}
 result:=GameEnding{Type:outcome,Accusation:a}
 switch outcome {
 case "win":result.Title="Дело раскрыто";result.Text=s.Case.Endings.Win
 case "fail":result.Title="Неверное обвинение";result.Text=s.Case.Endings.Fail
 default:result.Title="Недостаточно доказательств";result.Text=s.Case.Endings.Unsolved
 }
 // Rules customize the text of the outcome without changing the accusation verdict.
 bestPriority:=0
 matched:=false
 for _,rule:=range s.EndingRules {
  if rule.ID!=outcome || !RequirementsMet(s,rule.Requirements){continue}
  if !matched || rule.Priority>bestPriority {
   result.Title=rule.Title;result.Text=rule.Text;bestPriority=rule.Priority;matched=true
  }
 }
 s.Ending=&result
 if s.Clock!=nil{s.Clock.Pause()}
 return result,nil
}

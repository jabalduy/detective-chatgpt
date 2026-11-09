package game

import (
 "errors"
 "time"
)

// Hint is authored case content. Requirements prevent spoilers before
// relevant discoveries; CostSeconds overrides the default escalation.
type Hint struct {
 ID string `json:"id"`
 Text string `json:"text"`
 Requirements []Requirement `json:"requirements,omitempty"`
 CostSeconds int `json:"cost_seconds,omitempty"`
}

type HintView struct {
 ID string `json:"id"`
 Text string `json:"text"`
 CostSeconds int `json:"cost_seconds"`
 Number int `json:"number"`
}

// HintCostSeconds escalates for repeated hints and is capped at 15 minutes.
func HintCostSeconds(used int) int {
 costs:=[]int{180,300,480,720,900}
 if used<0 {used=0}
 if used>=len(costs){return costs[len(costs)-1]}
 return costs[used]
}

func AvailableHint(s *GameState) (HintView,error) {
 if !s.CanInvestigateAt(time.Now()) {return HintView{},errors.New("investigation time expired")}
 for _,hint:=range s.Hints {
  if s.UsedHints[hint.ID] || !RequirementsMet(s,hint.Requirements) {continue}
  cost:=hint.CostSeconds
  if cost<=0 {cost=HintCostSeconds(s.HintCount)}
  return HintView{ID:hint.ID,Text:hint.Text,CostSeconds:cost,Number:s.HintCount+1},nil
 }
 return HintView{},errors.New("no hint available")
}

func TakeHint(s *GameState) (HintView,error) {
 view,err:=AvailableHint(s)
 if err!=nil{return HintView{},err}
 if s.UsedHints==nil {s.UsedHints=make(map[string]bool)}
 s.UsedHints[view.ID]=true
 s.HintCount++
 if s.Clock!=nil {s.Clock.Add(time.Duration(view.CostSeconds)*time.Second)}
 return view,nil
}

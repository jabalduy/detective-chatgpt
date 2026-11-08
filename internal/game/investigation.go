package game

import (
 "fmt"
 "time"
)

// Examination definitions are authored in examinations.json.
// Completion is based on game time, allowing investigation to continue meanwhile.
type Examination struct {
 ID string `json:"id"`
 Title string `json:"title"`
 DurationSeconds int `json:"duration_seconds"`
 Requirements []Requirement `json:"requirements,omitempty"`
 Result string `json:"result"`
 Effects []Effect `json:"effects,omitempty"`
}
type ExaminationProgress struct {
 ID string `json:"id"`
 Title string `json:"title"`
 Status string `json:"status"`
 RemainingSeconds int `json:"remaining_seconds"`
 Result string `json:"result,omitempty"`
}
func examinationByID(s *GameState,id string)*Examination {
 for i:=range s.Examinations {if s.Examinations[i].ID==id{return &s.Examinations[i]}}
 return nil
}
func StartExamination(s *GameState,id string)error {
 if s.Ending!=nil {return fmt.Errorf("game finished")}
 e:=examinationByID(s,id)
 if e==nil {return fmt.Errorf("examination not found")}
 if s.Clock==nil {return fmt.Errorf("game clock unavailable")}
 if e.DurationSeconds<0 {return fmt.Errorf("invalid duration")}
 if !RequirementsMet(s,e.Requirements){return fmt.Errorf("required evidence missing")}
 if s.ExaminationStarted==nil{s.ExaminationStarted=make(map[string]time.Time)}
 if _,exists:=s.ExaminationStarted[id];exists{return fmt.Errorf("examination already started")}
 s.ExaminationStarted[id]=s.Clock.GetTime()
 return nil
}
func TickExaminations(s *GameState){
 if s.Clock==nil || s.Ending!=nil{return}
 if s.ExaminationCompleted==nil{s.ExaminationCompleted=make(map[string]bool)}
 for _,e:=range s.Examinations {
  started,ok:=s.ExaminationStarted[e.ID]
  if !ok || s.ExaminationCompleted[e.ID]{continue}
  if s.Clock.GetTime().Sub(started)<time.Duration(e.DurationSeconds)*time.Second {continue}
  s.ExaminationCompleted[e.ID]=true
  ApplyEffects(s,e.Effects)
  UpdateOpenDialogues(s)
 }
}
func ExaminationViews(s *GameState)[]ExaminationProgress{
 TickExaminations(s)
 views:=make([]ExaminationProgress,0,len(s.Examinations))
 for _,e:=range s.Examinations {
  v:=ExaminationProgress{ID:e.ID,Title:e.Title,Status:"locked"}
  if RequirementsMet(s,e.Requirements){v.Status="available"}
  if started,ok:=s.ExaminationStarted[e.ID];ok{
   v.Status="running"
   if s.Clock!=nil {
    remaining:=time.Duration(e.DurationSeconds)*time.Second-s.Clock.GetTime().Sub(started)
    if remaining>0 {v.RemainingSeconds=int((remaining+time.Second-1)/time.Second)}
   }
  }
  if s.ExaminationCompleted[e.ID]{v.Status="completed";v.RemainingSeconds=0;v.Result=e.Result}
  views=append(views,v)
 }
 return views
}

// EvidencePresentation lets the player challenge a character with a discovered clue/fact.
type EvidencePresentation struct {
 ID string `json:"id"`
 CharacterID int `json:"character_id"`
 EvidenceType string `json:"evidence_type"`
 EvidenceID int `json:"evidence_id"`
 Response string `json:"response"`
 Requirements []Requirement `json:"requirements,omitempty"`
 Effects []Effect `json:"effects,omitempty"`
 TimeCost int `json:"time_cost,omitempty"`
}
type EvidenceView struct {
 ID string `json:"id"`
 CharacterID int `json:"character_id"`
 EvidenceType string `json:"evidence_type"`
 EvidenceID int `json:"evidence_id"`
}
func evidenceAvailable(s *GameState,p EvidencePresentation)bool{
 if s.PresentedEvidence[p.ID] || !RequirementsMet(s,p.Requirements){return false}
 switch p.EvidenceType {
 case "clue":if !s.FoundClues[p.EvidenceID]{return false}
 case "fact":if !s.FoundFacts[p.EvidenceID]{return false}
 default:return false
 }
 character,ok:=s.Characters[p.CharacterID]
 return ok && character.Present && character.Location==s.CurrentScene
}
func AvailableEvidence(s *GameState)[]EvidenceView{
 views:=make([]EvidenceView,0)
 for _,p:=range s.EvidencePresentations {
  if evidenceAvailable(s,p){views=append(views,EvidenceView{p.ID,p.CharacterID,p.EvidenceType,p.EvidenceID})}
 }
 return views
}
func PresentEvidence(s *GameState,id string)(string,error){
 if s.Ending!=nil {return "",fmt.Errorf("game finished")}
 if s.ActiveEvent!=nil || s.ActiveConversation!="" || s.ActiveDialogue!=nil {return "",fmt.Errorf("finish active interaction first")}
 for _,p:=range s.EvidencePresentations {
  if p.ID!=id{continue}
  if !evidenceAvailable(s,p){return "",fmt.Errorf("evidence presentation unavailable")}
  if !PerformAction(s,Action{Effects:p.Effects,TimeCost:p.TimeCost}){return "",fmt.Errorf("action failed")}
  if s.PresentedEvidence==nil{s.PresentedEvidence=make(map[string]bool)}
  s.PresentedEvidence[id]=true
  return p.Response,nil
 }
 return "",fmt.Errorf("evidence presentation not found")
}

package game

import (
 "errors"
 "time"
)

var (
 ErrDiscoveryTooEarly = errors.New("official discovery is not yet due")
 ErrDiscoveryInteractionActive = errors.New("finish the active interaction before discovery")
)

// BeginInvestigation is the single, idempotent narrative transition.
// Call under the session mutex and persist the resulting state atomically.
// It does not itself advance the story clock or interrupt urgent choices.
func (s *GameState) BeginInvestigation(now time.Time) (bool,error) {
 if s.Ending!=nil {return false,errors.New("game already finished")}
 if s.CurrentNarrativeAct()!=ActPrologue {return false,nil}
 if s.Clock==nil {return false,ErrDiscoveryTooEarly}
 storyNow:=s.Clock.GetTime()
 discovery:=time.Date(storyNow.Year(),storyNow.Month(),storyNow.Day(),22,17,0,0,storyNow.Location())
 if storyNow.Before(discovery) {return false,ErrDiscoveryTooEarly}
 if s.ActiveEvent!=nil || s.ActiveConversation!="" || s.ActiveDialogue!=nil {
  return false,ErrDiscoveryInteractionActive
 }
 s.NarrativeAct=ActInvestigation
 s.DiscoveryCompleted=true
 if s.InvestigationDeadline.IsZero() {
  s.InvestigationDeadline=now.Add(60*time.Minute)
 }
 if s.AccusationDeadline.IsZero() {
  s.AccusationDeadline=s.InvestigationDeadline.Add(5*time.Minute)
 }
 return true,nil
}

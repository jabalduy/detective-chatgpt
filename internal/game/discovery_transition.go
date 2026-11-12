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
 // Discovery is anchored to the departure day, not the current calendar
 // day: a delayed prologue can continue past midnight.
 discoveryDay:=storyNow
 if !s.RouteStart.IsZero() {discoveryDay=s.RouteStart.In(storyNow.Location())}
 discovery:=time.Date(discoveryDay.Year(),discoveryDay.Month(),discoveryDay.Day(),22,17,0,0,storyNow.Location())
 if storyNow.Before(discovery) {return false,ErrDiscoveryTooEarly}
 if s.ActiveEvent!=nil || s.ActiveConversation!="" || s.ActiveDialogue!=nil {
  return false,ErrDiscoveryInteractionActive
 }
 s.NarrativeAct=ActInvestigation
 s.DiscoveryCompleted=true
 // The timer begins at official discovery, never at the start of the train ride.
 // A prologue may have stale deadlines from an older initialization path.
 s.InvestigationDeadline=now.Add(60*time.Minute)
 s.AccusationDeadline=s.InvestigationDeadline.Add(5*time.Minute)
 return true,nil
}

package game

import "time"

// RealTimePhase is based on wall-clock time, not the train clock.
// A conversation already in progress may finish after the investigation
// deadline, but it does not extend the accusation deadline.
type RealTimePhase string

const (
 PhaseInvestigation RealTimePhase = "investigation"
 PhaseConversationGrace RealTimePhase = "conversation_grace"
 PhaseAccusation RealTimePhase = "accusation"
 PhaseExpired RealTimePhase = "expired"
 PhaseFinished RealTimePhase = "finished"
)

func (s *GameState) PhaseAt(now time.Time) RealTimePhase {
 if s.Ending != nil { return PhaseFinished }
 // Legacy/test states without deadlines remain compatible until migrated.
 if s.InvestigationDeadline.IsZero() { return PhaseInvestigation }
 if now.Before(s.InvestigationDeadline) { return PhaseInvestigation }
 if !s.AccusationDeadline.IsZero() && !now.Before(s.AccusationDeadline) {
  return PhaseExpired
 }
 if s.ActiveConversation != "" || s.ActiveDialogue != nil {
  return PhaseConversationGrace
 }
 return PhaseAccusation
}

func (s *GameState) CanInvestigateAt(now time.Time) bool {
 return s.PhaseAt(now) == PhaseInvestigation
}

func (s *GameState) CanAccuseAt(now time.Time) bool {
 return s.PhaseAt(now) == PhaseAccusation
}

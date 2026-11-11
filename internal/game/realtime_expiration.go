package game

import "time"

// ExpireAccusation records a deterministic unresolved outcome once the
// accusation deadline passes. An already-started conversation may finish
// first; finishing after the original deadline grants five full minutes.
func (s *GameState) ExpireAccusation(now time.Time) bool {
 if s.Ending!=nil || s.AccusationDeadline.IsZero() || now.Before(s.AccusationDeadline) {
  return false
 }
 if s.ActiveConversation!="" || s.ActiveDialogue!=nil {return false}
 s.Ending=&GameEnding{
  Type:"unsolved",
  Title:"Время вышло",
  Text:s.Case.Endings.Unsolved,
 }
 if s.Clock!=nil {s.Clock.Pause()}
 return true
}

// FinishConversationGrace preserves the original 65-minute deadline when
// the conversation finishes early, or grants five full minutes when late.
// Call only after an already-active conversation/dialogue has ended.
func (s *GameState) FinishConversationGrace(now time.Time) bool {
 if s.Ending != nil || s.InvestigationDeadline.IsZero() || now.Before(s.InvestigationDeadline) {
  return false
 }
 if s.ActiveConversation != "" || s.ActiveDialogue != nil { return false }
 if s.AccusationDeadline.IsZero() || !now.Before(s.AccusationDeadline) {
  s.AccusationDeadline = now.Add(5 * time.Minute)
  return true
 }
 return false
}

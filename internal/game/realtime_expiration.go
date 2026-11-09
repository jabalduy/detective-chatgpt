package game

import "time"

// ExpireAccusation records a deterministic unresolved outcome once the
// accusation deadline passes. An already-started conversation may finish
// first; it cannot extend the fixed accusation deadline.
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

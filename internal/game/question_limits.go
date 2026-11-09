package game

import "fmt"

// CanAskDialogue enforces authored per-suspect question budgets. A zero
// limit means unlimited, preserving compatibility with older case files.
func CanAskDialogue(s *GameState, suspectID int) error {
 var limit int
 found:=false
 for _,suspect:=range s.Suspects {
  if suspect.SusID==suspectID {limit=suspect.MaxQuestions;found=true;break}
 }
 if !found{return fmt.Errorf("suspect not found")}
 if limit<=0{return nil}
 if s.AskedQuestionsBySuspect==nil{return nil}
 if s.AskedQuestionsBySuspect[suspectID]>=limit {
  return fmt.Errorf("question limit reached for suspect")
 }
 return nil
}

func RecordDialogueQuestion(s *GameState,suspectID int) {
 if s.AskedQuestionsBySuspect==nil {s.AskedQuestionsBySuspect=make(map[int]int)}
 s.AskedQuestionsBySuspect[suspectID]++
}

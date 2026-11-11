package game

import (
 "testing"
 "time"
)

func TestFinishConversationGrace(t *testing.T) {
 start:=time.Date(2026,11,11,12,0,0,0,time.UTC)
 for _,tc:=range []struct{name string; end time.Duration; want time.Duration; extended bool}{
  {"before investigation end",59*time.Minute,65*time.Minute,false},
  {"early grace",62*time.Minute,65*time.Minute,false},
  {"at accusation deadline",65*time.Minute,70*time.Minute,true},
  {"late grace",70*time.Minute,75*time.Minute,true},
 }{
  t.Run(tc.name,func(t *testing.T){
   s:=&GameState{InvestigationDeadline:start.Add(60*time.Minute),AccusationDeadline:start.Add(65*time.Minute),ActiveConversation:"witness"}
   now:=start.Add(tc.end)
   if s.FinishConversationGrace(now){t.Fatal("cannot extend while conversation is active")}
   s.ActiveConversation=""
   if got:=s.FinishConversationGrace(now);got!=tc.extended{t.Fatalf("extended=%v want %v",got,tc.extended)}
   if want:=start.Add(tc.want);!s.AccusationDeadline.Equal(want){t.Fatalf("deadline=%v want %v",s.AccusationDeadline,want)}
   if tc.end>=60*time.Minute && s.PhaseAt(now)!=PhaseAccusation{t.Fatal("accusation not available")}
  })
 }
}

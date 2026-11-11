package game

import (
 "testing"
 "time"
)

func TestInvestigationDeadlineConversationGraceAndAccusationWindow(t *testing.T) {
 start:=time.Date(1928,10,1,21,0,0,0,time.UTC)
 s:=&GameState{
  Case:CaseInfo{Endings:Endings{Unsolved:"The case remains open."}},
  InvestigationDeadline:start.Add(60*time.Minute),
  AccusationDeadline:start.Add(65*time.Minute),
  CurrentScene:"salon",
  Characters:map[int]CharacterState{2:{Location:"salon",Present:true}},
  Conversations:[]Conversation{{ID:"boris",CharacterID:2,StartNode:"hello",Nodes:[]ConversationNode{{
   ID:"hello",Choices:[]ConversationChoice{{ID:"leave",End:true}},
  }}}},
 }
 if err:=StartConversation(s,"boris");err!=nil {t.Fatal(err)}
 late:=start.Add(67*time.Minute)
 if got:=s.PhaseAt(late);got!=PhaseConversationGrace {t.Fatalf("phase during active dialogue = %s",got)}
 if s.ExpireAccusation(late) {t.Fatal("active dialogue expired before it could finish")}
 if _,err:=ChooseConversationOption(s,"leave");err!=nil {t.Fatal(err)}
 if !s.FinishConversationGrace(late) {t.Fatal("late dialogue should extend accusation deadline")}
 if want:=late.Add(5*time.Minute);!s.AccusationDeadline.Equal(want) {
  t.Fatalf("deadline = %s want %s",s.AccusationDeadline,want)
 }
 if !s.CanAccuseAt(late) {t.Fatal("accusation unavailable after grace")}
 if s.CanInvestigateAt(late) {t.Fatal("investigation reopened after grace")}
 if s.ExpireAccusation(late.Add(5*time.Minute-time.Nanosecond)) {t.Fatal("accusation expired early")}
 if !s.ExpireAccusation(late.Add(5*time.Minute)) {t.Fatal("accusation did not expire at deadline")}
 if s.Ending==nil || s.Ending.Type!="unsolved" {t.Fatal("missing unresolved timeout ending")}
 if s.ExpireAccusation(late.Add(6*time.Minute)) {t.Fatal("timeout ending recorded twice")}
 if s.PhaseAt(late)!=PhaseFinished {t.Fatal("finished phase not preserved")}
}

func TestConversationGraceDoesNotShortenExistingAccusationWindow(t *testing.T){
 start:=time.Date(1928,10,1,21,0,0,0,time.UTC)
 s:=&GameState{InvestigationDeadline:start.Add(time.Hour),AccusationDeadline:start.Add(65*time.Minute)}
 s.ActiveConversation="boris"
 at:=start.Add(61*time.Minute)
 if s.ExpireAccusation(at) {t.Fatal("active dialogue expired")}
 s.ActiveConversation=""
 if s.FinishConversationGrace(at) {t.Fatal("existing deadline should be preserved")}
 if want:=start.Add(65*time.Minute);!s.AccusationDeadline.Equal(want) {t.Fatalf("deadline changed to %v",s.AccusationDeadline)}
}

func TestEndingFreezesGameAndRejectsSecondAccusation(t *testing.T){
 s:=endingTestState()
 defer s.Clock.StopTime()
 s.InvestigationDeadline=time.Now().Add(-time.Hour)
 s.AccusationDeadline=time.Now().Add(time.Minute)
 ending,err:=SubmitAccusation(s,Accusation{SuspectID:1,MotiveID:2,FactIDs:[]int{3,4}})
 if err!=nil || ending.Type!="win" {t.Fatalf("ending=%+v err=%v",ending,err)}
 if !s.Clock.IsPaused() {t.Fatal("route clock still running after verdict")}
 if s.PhaseAt(time.Now())!=PhaseFinished {t.Fatal("phase not finished")}
 if _,err:=SubmitAccusation(s,Accusation{SuspectID:1,MotiveID:2,FactIDs:[]int{3,4}});err==nil {t.Fatal("second accusation accepted")}
 if s.ExpireAccusation(time.Now().Add(2*time.Hour)) {t.Fatal("timeout overwrote verdict")}
 if PerformAction(s,Action{}) {t.Fatal("action permitted after ending")}
 if TriggerEvent(s,123) {t.Fatal("event triggered after ending")}
}

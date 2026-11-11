package game

import (
 "testing"
 "time"
)

// This synthetic case proves that a missed timed decision does not make the
// only winning evidence path unreachable. It is not the authored story.
func TestAlternateEvidencePathsAfterUrgentDecision(t *testing.T) {
 for _,tc:=range []struct{name string;choose bool}{
  {"successful urgent decision",true},
  {"missed urgent decision",false},
 } {
  t.Run(tc.name,func(t *testing.T){
   start:=time.Date(2026,1,1,22,0,0,0,time.UTC)
   s:=&GameState{
    Clock:NewClock(start),CurrentScene:"corridor",
    Case:CaseInfo{Solution:Solution{SuspectID:1,MotiveID:2,FactIDs:[]int{30}},Endings:Endings{Win:"Solved"}},
    Motives:[]Motive{{ID:2,SuspectID:1}},
    FoundClues:map[int]bool{},FoundFacts:map[int]bool{},FoundMotives:map[int]bool{2:true},
    Events:[]Event{{ID:5,TimeLimit:10,DefaultOptionID:20,Options:[]EventOption{
     {ID:10,Text:"Catch witness",Action:Action{Effects:[]Effect{{Type:"discover_fact",ID:30}}}},
     {ID:20,Text:"Witness escapes",Action:Action{Effects:[]Effect{{Type:"discover_fact",ID:40}}}},
    }}},
    Scenes:[]Scene{{ID:"corridor",Choices:[]Choice{{ID:"inspect",Action:Action{Effects:[]Effect{{Type:"discover_clue",ID:11}}}}}}},
    Deductions:[]Deduction{{ID:3,Required:[]int{40},Action:Action{Requirements:[]Requirement{{Type:"clue_found",ID:11}},Effects:[]Effect{{Type:"discover_fact",ID:30}}}}},
   }
   defer s.Clock.StopTime()
   if !TriggerEvent(s,5){t.Fatal("event not triggered")}
   if tc.choose {
    if !ChooseEventOption(s,10){t.Fatal("urgent choice failed")}
   } else {
    s.EventDeadline=time.Now().Add(-time.Second)
    if !ExpireActiveEvent(s){t.Fatal("missed event did not resolve")}
    if !s.FoundFacts[40]{t.Fatal("default consequence missing")}
    if err:=ChooseSceneOption(s,"inspect");err!=nil{t.Fatal(err)}
    if !SolveDeduction(s,s.Deductions[0]){t.Fatal("alternate deduction path blocked")}
   }
   if !s.FoundFacts[30]{t.Fatal("required evidence unreachable")}
   ending,err:=SubmitAccusation(s,Accusation{SuspectID:1,MotiveID:2,FactIDs:[]int{30}})
   if err!=nil || ending.Type!="win"{t.Fatalf("winning path failed: %+v %v",ending,err)}
  })
 }
}

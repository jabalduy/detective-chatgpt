package game

import (
 "testing"
 "time"
)

func TestRepeatedSceneActionIsBlockedByRequirements(t *testing.T){
 s:=&GameState{
  CurrentScene:"corridor",
  Scenes:[]Scene{{ID:"corridor",Choices:[]Choice{{
   ID:"search",Action:Action{
    Requirements:[]Requirement{{Type:"clue_found",ID:99}},
    Effects:[]Effect{{Type:"discover_fact",ID:12}},
   },
  }}}},
  FoundClues:map[int]bool{},FoundFacts:map[int]bool{},
 }
 if err:=ChooseSceneOption(s,"search");err==nil{t.Fatal("locked action accepted")}
 if s.FoundFacts[12]{t.Fatal("locked action granted evidence")}
 s.FoundClues[99]=true
 if err:=ChooseSceneOption(s,"search");err!=nil{t.Fatal(err)}
 if !s.FoundFacts[12]{t.Fatal("unlocked action did not grant fact")}
}

func TestExpiredEventCannotApplyLateChoiceTwice(t *testing.T){
 s:=&GameState{
  FoundFacts:map[int]bool{},
  Events:[]Event{{ID:1,TimeLimit:5,DefaultOptionID:20,Options:[]EventOption{
   {ID:10,Action:Action{Effects:[]Effect{{Type:"discover_fact",ID:10}}}},
   {ID:20,Action:Action{Effects:[]Effect{{Type:"discover_fact",ID:20}}}},
  }}},
 }
 if !TriggerEvent(s,1){t.Fatal("event not triggered")}
 s.EventDeadline=time.Now().Add(-time.Second)
 if ChooseEventOption(s,10){t.Fatal("late option accepted")}
 if !s.FoundFacts[20] || s.FoundFacts[10]{t.Fatal("default effect missing or late effect applied")}
 if ChooseEventOption(s,10) || ExpireActiveEvent(s){t.Fatal("resolved event accepted a second action")}
}

func TestGameEndingBlocksAllFurtherActions(t *testing.T){
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 s:=&GameState{
  Clock:NewClock(start),CurrentScene:"corridor",
  Scenes:[]Scene{{ID:"corridor",Choices:[]Choice{{ID:"search",Action:Action{Effects:[]Effect{{Type:"discover_fact",ID:99}}}}}}},
  FoundFacts:map[int]bool{1:true},FoundMotives:map[int]bool{2:true},
  Motives:[]Motive{{ID:2,SuspectID:7}},
  Case:CaseInfo{Solution:Solution{SuspectID:7,MotiveID:2,FactIDs:[]int{1}}},
  Events:[]Event{{ID:5}},
 }
 defer s.Clock.StopTime()
 result,err:=SubmitAccusation(s,Accusation{SuspectID:7,MotiveID:2,FactIDs:[]int{1}})
 if err!=nil || result.Type!="win"{t.Fatalf("accusation: %+v %v",result,err)}
 if !s.Clock.IsPaused(){t.Fatal("ending did not pause clock")}
 if err:=ChooseSceneOption(s,"search");err==nil{t.Fatal("scene action accepted after ending")}
 if err:=StartConversation(s,"missing");err==nil{t.Fatal("conversation accepted after ending")}
 if TriggerEvent(s,5){t.Fatal("event triggered after ending")}
 if _,err:=SubmitAccusation(s,Accusation{SuspectID:7,MotiveID:2,FactIDs:[]int{1}});err==nil{
  t.Fatal("second accusation accepted")
 }
 if s.FoundFacts[99]{t.Fatal("post-ending action changed facts")}
}

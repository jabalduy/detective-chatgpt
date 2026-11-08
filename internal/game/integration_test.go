package game

import (
 "testing"
 "time"
)

// TestFullInvestigationFlow exercises the engine as a single playable sequence.
// It deliberately uses a tiny synthetic case, not the final murder mystery.
func TestFullInvestigationFlow(t *testing.T) {
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 s:=&GameState{
  Clock:NewClock(start),RouteStart:start,CurrentScene:"corridor",
  Scenes:[]Scene{
   {ID:"corridor",Choices:[]Choice{{ID:"inspect",Action:Action{TimeCost:30,Effects:[]Effect{{Type:"discover_clue",ID:10}}}},{ID:"move",NextScene:"compartment",Action:Action{TimeCost:20}}}},
   {ID:"compartment",Choices:[]Choice{{ID:"return",NextScene:"corridor"}}},
  },
  Characters:map[int]CharacterState{1:{Location:"corridor",Present:true}},
  Conversations:[]Conversation{{ID:"witness",CharacterID:1,StartNode:"first",Nodes:[]ConversationNode{
   {ID:"first",Text:"What happened?",Choices:[]ConversationChoice{{ID:"ask",NextNode:"last",Action:Action{Effects:[]Effect{{Type:"discover_motive",ID:2}}}}}},
   {ID:"last",Text:"I saw something.",Choices:[]ConversationChoice{{ID:"bye",End:true}}},
  }}},
  Examinations:[]Examination{{ID:"sample",DurationSeconds:120,Requirements:[]Requirement{{Type:"clue_found",ID:10}},Result:"Match",Effects:[]Effect{{Type:"discover_fact",ID:20}}}},
  EvidencePresentations:[]EvidencePresentation{{ID:"show",CharacterID:1,EvidenceType:"fact",EvidenceID:20,Response:"I admit it.",Effects:[]Effect{{Type:"discover_fact",ID:21}}}},
  Stations:[]Station{{ID:"stop",Name:"Stop",Arrival:"00:00",Departure:"00:02"}},
  Events:[]Event{{ID:99,At:"00:01",TimeLimit:10,DefaultOptionID:1,Options:[]EventOption{{ID:1,Text:"Wait"},{ID:2,Text:"Act"}}}},
  Case:CaseInfo{Solution:Solution{SuspectID:1,MotiveID:2,FactIDs:[]int{20,21}},Endings:Endings{Win:"Solved",Fail:"Wrong",Unsolved:"Incomplete"}},
  Motives:[]Motive{{ID:2,SuspectID:1}},
  FoundClues:make(map[int]bool),FoundFacts:make(map[int]bool),FoundMotives:make(map[int]bool),
  OpenDialogues:make(map[DialogueID]bool),
 }
 if err:=ChooseSceneOption(s,"inspect");err!=nil{t.Fatal(err)}
 if !s.FoundClues[10]{t.Fatal("clue missing")}
 if err:=StartExamination(s,"sample");err!=nil{t.Fatal(err)}
 if err:=StartConversation(s,"witness");err!=nil{t.Fatal(err)}
 if err:=ChooseSceneOption(s,"move");err==nil{t.Fatal("moved during conversation")}
 if _,err:=ChooseConversationOption(s,"ask");err!=nil{t.Fatal(err)}
 ended,err:=ChooseConversationOption(s,"bye")
 if err!=nil||!ended||!s.FoundMotives[2]{t.Fatalf("conversation end=%v err=%v",ended,err)}
 if err:=ChooseSceneOption(s,"move");err!=nil{t.Fatal(err)}
 if s.CurrentScene!="compartment"{t.Fatal("movement failed")}
 if err:=ChooseSceneOption(s,"return");err!=nil{t.Fatal(err)}
 s.Clock.Add(120*time.Second)
 views:=ExaminationViews(s)
 if len(views)!=1||views[0].Status!="completed"||!s.FoundFacts[20]{t.Fatalf("forensics: %+v",views)}
 if _,err:=PresentEvidence(s,"show");err!=nil{t.Fatal(err)}
 if !s.FoundFacts[21]{t.Fatal("confrontation effect missing")}
 s.Clock.Add(9*time.Minute)
 TickTrain(s)
 TickEvents(s)
 if s.ActiveEvent==nil {t.Fatal("scheduled event not activated")}
 if err:=ChooseSceneOption(s,"move");err==nil{t.Fatal("moved during urgent event")}
 if !ChooseEventOption(s,2){t.Fatal("urgent decision failed")}
 if s.ActiveEvent!=nil{t.Fatal("event still active")}
 ending,err:=SubmitAccusation(s,Accusation{SuspectID:1,MotiveID:2,FactIDs:[]int{20,21}})
 if err!=nil||ending.Type!="win"{t.Fatalf("ending=%+v err=%v",ending,err)}
 if err:=ChooseSceneOption(s,"move");err==nil{t.Fatal("game action allowed after ending")}
}

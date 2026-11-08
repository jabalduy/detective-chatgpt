package game

import (
 "testing"
 "time"
)

func TestConversationBranchesAndResumesClock(t *testing.T){
 start:=time.Date(2026,1,1,22,0,0,0,time.UTC)
 s:=&GameState{Clock:NewClock(start),CurrentScene:"corridor",
  Characters:map[int]CharacterState{1:{Location:"corridor",Present:true}},
  Conversations:[]Conversation{{ID:"intro",CharacterID:1,StartNode:"a",Nodes:[]ConversationNode{
   {ID:"a",Text:"Hello",Choices:[]ConversationChoice{{ID:"reply",NextNode:"b",Action:Action{TimeCost:30,Effects:[]Effect{{Type:"character_trust_up",ID:1}}}}}},
   {ID:"b",Text:"Thanks",Choices:[]ConversationChoice{{ID:"bye",End:true}}},
  }}},
 }
 if err:=StartConversation(s,"intro");err!=nil{t.Fatal(err)}
 view,err:=CurrentConversationView(s);if err!=nil||len(view.Choices)!=1{t.Fatalf("view=%+v err=%v",view,err)}
 ended,err:=ChooseConversationOption(s,"reply");if err!=nil||ended{t.Fatalf("ended=%v err=%v",ended,err)}
 if s.Characters[1].Trust!=1{t.Fatal("trust effect missing")}
 if s.Clock.GetTime().Sub(start)!=30*time.Second{t.Fatal("time cost missing")}
 ended,err=ChooseConversationOption(s,"bye");if err!=nil||!ended{t.Fatalf("ended=%v err=%v",ended,err)}
 if s.ActiveConversation!=""{t.Fatal("conversation not closed")}
}

func TestConversationRejectsAbsentCharacter(t *testing.T){
 s:=&GameState{CurrentScene:"corridor",Characters:map[int]CharacterState{1:{Location:"corridor",Present:false}},
 Conversations:[]Conversation{{ID:"intro",CharacterID:1,StartNode:"a",Nodes:[]ConversationNode{{ID:"a"}}}}}
 if err:=StartConversation(s,"intro");err==nil{t.Fatal("expected unavailable character")}
}

package game

import "testing"

func TestTrustUnlocksConversationChoice(t *testing.T){
 s:=&GameState{
  CurrentScene:"carriage",
  Characters:map[int]CharacterState{7:{Location:"carriage",Present:true,Trust:0}},
  Conversations:[]Conversation{{ID:"witness",CharacterID:7,StartNode:"intro",Nodes:[]ConversationNode{
   {ID:"intro",Choices:[]ConversationChoice{
    {ID:"secret",End:true,Action:Action{Requirements:[]Requirement{{Type:"character_trust_at_least",ID:7,MinTrust:2}},Effects:[]Effect{{Type:"discover_fact",ID:50}}}},
    {ID:"kind",NextNode:"intro",Action:Action{Effects:[]Effect{{Type:"character_trust_up",ID:7}}}},
    {ID:"insult",NextNode:"intro",Action:Action{Effects:[]Effect{{Type:"character_trust_down",ID:7}}}},
   }},
  }}},
 }
 if err:=StartConversation(s,"witness");err!=nil{t.Fatal(err)}
 view,err:=CurrentConversationView(s)
 if err!=nil{t.Fatal(err)}
 for _,choice:=range view.Choices{if choice.ID=="secret"{t.Fatal("secret should be hidden")}}
 if _,err:=ChooseConversationOption(s,"secret");err==nil{t.Fatal("hidden choice should be rejected")}
 if s.FoundFacts[50]{t.Fatal("rejected action changed game state")}
 for i:=0;i<2;i++ {if _,err:=ChooseConversationOption(s,"kind");err!=nil{t.Fatal(err)}}
 view,err=CurrentConversationView(s);if err!=nil{t.Fatal(err)}
 found:=false
 for _,choice:=range view.Choices{if choice.ID=="secret"{found=true}}
 if !found{t.Fatal("trust did not unlock secret")}
 if _,err:=ChooseConversationOption(s,"insult");err!=nil{t.Fatal(err)}
 if _,err:=ChooseConversationOption(s,"secret");err==nil{t.Fatal("lost trust should relock secret")}
 if _,err:=ChooseConversationOption(s,"kind");err!=nil{t.Fatal(err)}
 ended,err:=ChooseConversationOption(s,"secret")
 if err!=nil || !ended || !s.FoundFacts[50]{t.Fatalf("secret not revealed: ended=%v err=%v",ended,err)}
}

func TestTrustRequirementRejectsMissingCharacter(t *testing.T){
 s:=&GameState{}
 if RequirementsMet(s,[]Requirement{{Type:"character_trust_at_least",ID:999,MinTrust:0}}){
  t.Fatal("missing character must not pass a zero threshold")
 }
}

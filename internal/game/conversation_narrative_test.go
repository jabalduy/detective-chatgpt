package game

import "testing"

func TestConversationNarrativeAvailability(t *testing.T) {
 for _,tc:=range []struct{name string;stateAct,conversationAct NarrativeAct;want bool}{
  {"prologue introduction",ActPrologue,ActPrologue,true},
  {"investigation hidden in prologue",ActPrologue,ActInvestigation,false},
  {"unclassified hidden in prologue",ActPrologue,"",false},
  {"introduction hidden after discovery",ActInvestigation,ActPrologue,false},
  {"legacy investigation","","",true},
  {"normal investigation",ActInvestigation,ActInvestigation,true},
 } {
  t.Run(tc.name,func(t *testing.T){
   s:=&GameState{NarrativeAct:tc.stateAct}
   c:=&Conversation{ID:"test",AvailableAct:tc.conversationAct}
   if got:=ConversationAvailableInAct(s,c);got!=tc.want {t.Fatalf("available=%v, want %v",got,tc.want)}
  })
 }
}

func TestStartConversationRejectsInvestigationDialogueInPrologue(t *testing.T) {
 s:=&GameState{NarrativeAct:ActPrologue,Conversations:[]Conversation{{ID:"investigate",CharacterID:1,StartNode:"start",Nodes:[]ConversationNode{{ID:"start"}}}},Characters:map[int]CharacterState{1:{Present:true,Location:"salon"}},CurrentScene:"salon"}
 if err:=StartConversation(s,"investigate");err==nil {t.Fatal("unclassified investigation conversation allowed during prologue")}
 if s.ActiveConversation!="" {t.Fatal("rejected conversation changed state")}
 s.Conversations[0].AvailableAct=ActPrologue
 if err:=StartConversation(s,"investigate");err!=nil {t.Fatalf("prologue conversation rejected: %v",err)}
}

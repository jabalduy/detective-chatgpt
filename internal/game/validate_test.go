package game

import (
 "strings"
 "testing"
)

func TestValidateNarrativeValid(t *testing.T){
 scenes:=[]Scene{
  {ID:"hall",Choices:[]Choice{{ID:"enter",NextScene:"room"}}},
  {ID:"room",Choices:[]Choice{{ID:"back",NextScene:"hall"}}},
 }
 conversations:=[]Conversation{{ID:"witness",StartNode:"a",Nodes:[]ConversationNode{
  {ID:"a",Choices:[]ConversationChoice{{ID:"ask",NextNode:"b"}}},
  {ID:"b",Choices:[]ConversationChoice{{ID:"repeat",NextNode:"a"},{ID:"leave",End:true}}},
 }}}
 if issues:=ValidateNarrative(scenes,conversations);len(issues)!=0{t.Fatalf("unexpected issues: %v",issues)}
}

func TestValidateNarrativeBrokenReferences(t *testing.T){
 scenes:=[]Scene{{ID:"hall",Choices:[]Choice{{ID:"enter",NextScene:"missing"},{ID:"enter"}}}}
 conversations:=[]Conversation{{ID:"witness",StartNode:"a",Nodes:[]ConversationNode{
  {ID:"a",Choices:[]ConversationChoice{{ID:"next",NextNode:"missing"}}},
 }}}
 issues:=ValidateNarrative(scenes,conversations)
 combined:=""
 for _,issue:=range issues {combined+=issue.Error()+"\n"}
 for _,expected:=range []string{"duplicate choice","missing scene","missing node","no path to an ending"}{
  if !strings.Contains(combined,expected){t.Fatalf("missing diagnostic %q in %s",expected,combined)}
 }
}

func TestValidateNarrativeDeadEndCycle(t *testing.T){
 conversations:=[]Conversation{{ID:"loop",StartNode:"a",Nodes:[]ConversationNode{
  {ID:"a",Choices:[]ConversationChoice{{ID:"forward",NextNode:"b"}}},
  {ID:"b",Choices:[]ConversationChoice{{ID:"back",NextNode:"a"}}},
 }}}
 issues:=ValidateNarrative(nil,conversations)
 if len(issues)!=2 {t.Fatalf("expected both nodes to be flagged: %v",issues)}
}

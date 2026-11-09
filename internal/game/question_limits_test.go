package game

import "testing"

func TestPerSuspectQuestionBudget(t *testing.T){
 s:=&GameState{Suspects:[]Suspect{{SusID:7,MaxQuestions:2}}}
 if err:=CanAskDialogue(s,7);err!=nil{t.Fatal(err)}
 RecordDialogueQuestion(s,7)
 if err:=CanAskDialogue(s,7);err!=nil{t.Fatal(err)}
 RecordDialogueQuestion(s,7)
 if err:=CanAskDialogue(s,7);err==nil{t.Fatal("third question should be blocked")}
 if err:=CanAskDialogue(s,8);err==nil{t.Fatal("unknown suspect should be rejected")}
 legacy:=&GameState{Suspects:[]Suspect{{SusID:9}}}
 for i:=0;i<5;i++ {if err:=CanAskDialogue(legacy,9);err!=nil{t.Fatal(err)}}
}

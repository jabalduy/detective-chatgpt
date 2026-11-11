package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

// The finale must be earned through the actual case data, not by setting
// solution facts or motives directly in the test state.
func TestExpress1928EndToEndVerdicts(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,issue:=range ValidateCaseIntegrity(def) {t.Error(issue)}
 for _,tc:=range []struct{
  name string
  suspect,motive int
  includeClockProof bool
  want string
 }{
  {"complete solution",1,1,true,"win"},
  {"correct killer insufficient proof",1,1,false,"unsolved"},
  {"plausible but wrong Boris accusation",2,2,false,"fail"},
 } {
  t.Run(tc.name,func(t *testing.T){
   s:=&game.GameState{
    Case:def.CaseInfo,Motives:def.Motives,EndingRules:def.EndingRules,
    Scenes:def.Scenes,Conversations:def.Conversations,Characters:def.Characters,
    Deductions:def.Deductions,Events:def.Events,CurrentScene:"corridor",
    FoundFacts:map[int]bool{},FoundMotives:map[int]bool{},
    SolvedDeductions:map[int]bool{},
   }
   step:=func(id string){
    t.Helper()
    if err:=game.ChooseSceneOption(s,id);err!=nil {t.Fatalf("scene %s choice %s: %v",s.CurrentScene,id,err)}
   }
   for _,id:=range []string{
    "victor","ribbon","paperweight","watch","note","back",
    "staff","box","nikola","cancellation","back",
    "salon","boris","rosa","back",
    "archive","fraud","family",
   } {step(id)}
   solve:=func(id int){
    t.Helper()
    for _,d:=range s.Deductions {
     if d.ID==id {
      if !game.SolveDeduction(s,d) {t.Fatalf("deduction %d unavailable",id)}
      return
     }
    }
    t.Fatalf("deduction %d missing",id)
   }
   for _,id:=range []int{1,2,4,5,7} {solve(id)}
   if tc.includeClockProof {
    step("back")
    step("salon")
    if err:=game.StartConversation(s,"boris_intro");err!=nil {t.Fatal(err)}
    if _,err:=game.ChooseConversationOption(s,"clock_confront");err!=nil {t.Fatal(err)}
    if _,err:=game.ChooseConversationOption(s,"record");err!=nil {t.Fatal(err)}
    solve(6)
   }
   facts:=[]int{10,11}
   if tc.includeClockProof {facts=append(facts,12)}
   if tc.suspect==2 {facts=[]int{4,13}}
   for _,id:=range facts {if !s.FoundFacts[id] {t.Fatalf("accusation fact %d not earned",id)}}
   if !s.FoundMotives[tc.motive] {t.Fatalf("motive %d not earned",tc.motive)}
   ending,err:=game.SubmitAccusation(s,game.Accusation{SuspectID:tc.suspect,MotiveID:tc.motive,FactIDs:facts})
   if err!=nil {t.Fatal(err)}
   if ending.Type!=tc.want {t.Fatalf("verdict = %s, want %s",ending.Type,tc.want)}
   if s.Ending==nil {t.Fatal("verdict not persisted")}
   if _,err:=game.SubmitAccusation(s,game.Accusation{SuspectID:tc.suspect,MotiveID:tc.motive,FactIDs:facts});err==nil {
    t.Fatal("second accusation accepted")
   }
  })
 }
}

func TestExpress1928MandatoryProofNotAwardedByScenesOrEvents(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 mandatory:=map[int]bool{}
 for _,id:=range def.Solution.FactIDs {mandatory[id]=true}
 check:=func(where string,effects []game.Effect){
  t.Helper()
  for _,e:=range effects {
   if e.Type=="discover_fact" && mandatory[e.ID] {t.Errorf("%s directly awards mandatory fact %d",where,e.ID)}
  }
 }
 for _,scene:=range def.Scenes {
  for _,choice:=range scene.Choices {check("scene "+scene.ID+" choice "+choice.ID,choice.Effects)}
 }
 for _,event:=range def.Events {
  for _,option:=range event.Options {check("event",option.Effects)}
 }
 for _,conversation:=range def.Conversations {
  for _,node:=range conversation.Nodes {
   for _,choice:=range node.Choices {check("conversation "+conversation.ID+" choice "+choice.ID,choice.Effects)}
  }
 }
 for _,d:=range def.Deductions {
  for _,e:=range d.Effects {
   if e.Type!="discover_fact" || !mandatory[e.ID] {continue}
   if len(d.Required)<2 {t.Errorf("mandatory fact %d has deduction with fewer than two prerequisites",e.ID)}
  }
 }
}

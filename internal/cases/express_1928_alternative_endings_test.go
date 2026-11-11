package cases

import (
 "path/filepath"
 "testing"
 "detective/internal/game"
)

func TestExpress1928AlternativeAccusations(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,tc:=range []struct{
  name string
  facts []int
  suspect int
  want string
 }{
  {"all mandatory evidence",[]int{10,11,12},1,"win"},
  {"missed clock proof",[]int{10,11},1,"unsolved"},
  {"missed weapon proof",[]int{10,12},1,"unsolved"},
  {"wrong suspect with true motive",[]int{10,11,12},2,"fail"},
 }{
  t.Run(tc.name,func(t *testing.T){
   s:=&game.GameState{Case:def.CaseInfo,Motives:def.Motives,EndingRules:def.EndingRules,
    FoundFacts:map[int]bool{10:true,11:true,12:true,16:true},
    FoundMotives:map[int]bool{1:true},
   }
   motive:=1
   if tc.suspect==2 {
    // Boris's motive is not discovered, so an accusation against him
    // must be rejected before a verdict is recorded.
    _,err:=game.SubmitAccusation(s,game.Accusation{SuspectID:2,MotiveID:motive,FactIDs:tc.facts})
    if err==nil || s.Ending!=nil {t.Fatal("invalid suspect/motive pairing must be rejected")}
    return
   }
   ending,err:=game.SubmitAccusation(s,game.Accusation{SuspectID:tc.suspect,MotiveID:motive,FactIDs:tc.facts})
   if err!=nil {t.Fatal(err)}
   if ending.Type!=tc.want {t.Fatalf("got %q want %q",ending.Type,tc.want)}
   if ending.Title=="" || ending.Text=="" {t.Fatal("ending narrative missing")}
  })
 }
}

func TestExpress1928OptionalEventChoicesDoNotGrantMandatoryEvidence(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 required:=map[int]bool{}
 for _,id:=range def.Solution.FactIDs {required[id]=true}
 for _,event:=range def.Events {
  for _,option:=range event.Options {
   for _,effect:=range option.Effects {
    if effect.Type=="discover_fact" && required[effect.ID] {
     t.Errorf("event %d option %d grants mandatory proof %d",event.ID,option.ID,effect.ID)
    }
   }
  }
 }
}

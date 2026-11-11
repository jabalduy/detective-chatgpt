package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

func TestExpress1928TimedEvents(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 if len(def.Events)!=4 {t.Fatalf("want four timed events, got %d",len(def.Events))}
 if problems:=game.ValidateEvents(def.Events);len(problems)>0 {t.Fatalf("invalid timed events: %v",problems)}
 factIDs:=map[int]bool{}
 for _,f:=range def.Facts {factIDs[f.ID]=true}
 ids:=map[int]bool{}
 for _,event:=range def.Events {
  if ids[event.ID] {t.Errorf("duplicate event %d",event.ID)}
  ids[event.ID]=true
  if event.TimeLimit!=10 {t.Errorf("event %d timer must be 10s",event.ID)}
  if event.At=="" {t.Errorf("event %d must have story clock trigger",event.ID)}
  hasDefault:=false
  for _,option:=range event.Options {
   if option.ID==event.DefaultOptionID {
    hasDefault=true
    if len(option.Requirements)>0 {t.Errorf("event %d default must be unconditional",event.ID)}
   }
   for _,effect:=range option.Effects {
    if effect.Type=="discover_fact" && !factIDs[effect.ID] {t.Errorf("event %d references missing fact %d",event.ID,effect.ID)}
   }
  }
  if !hasDefault {t.Errorf("event %d missing timeout option",event.ID)}
 }
 for _,id:=range def.Solution.FactIDs {
  if id==17 || id==18 || id==19 {t.Errorf("optional event fact %d must not be mandatory",id)}
 }
}

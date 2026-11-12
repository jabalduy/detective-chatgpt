package cases

import (
 "path/filepath"
 "testing"
 "time"
)

// Story-clock triggers must agree with the canonical 1928 timeline.
// The follow-up rooftop decision is triggered by the opening choice,
// never independently by the route clock.
func TestExpress1928CanonicalEventChronology(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 if got:=def.StartTime.StartTime.Format("15:04");got!="21:00" {
  t.Fatalf("case start=%s, want 21:00",got)
 }
 expected:=map[int]string{101:"21:12",102:"21:24",104:"21:51",103:"22:09",105:""}
 byID:=make(map[int]string,len(def.Events))
 for _,event:=range def.Events {
  byID[event.ID]=event.At
  want,ok:=expected[event.ID]
  if !ok {t.Errorf("unexpected event %d",event.ID);continue}
  if event.At!=want {t.Errorf("event %d at %q, want %q",event.ID,event.At,want)}
 }
 if len(byID)!=len(expected) {t.Fatalf("events=%v, want %d",byID,len(expected))}
 tunnel,err:=time.Parse("15:04","22:10")
 if err!=nil {t.Fatal(err)}
 chase,err:=time.Parse("15:04",byID[103])
 if err!=nil {t.Fatal(err)}
 if tunnel.Sub(chase)!=time.Minute {t.Fatalf("rooftop chase must precede tunnel by 1 minute, got %v",tunnel.Sub(chase))}
 var opening *int
 for _,event:=range def.Events {
  if event.ID!=103 {continue}
  for _,option:=range event.Options {
   if option.ID!=1 {continue}
   for _,effect:=range option.Effects {
    if effect.Type=="trigger_event" {id:=effect.ID;opening=&id}
   }
  }
 }
 if opening==nil || *opening!=105 {t.Fatal("roof event 105 must be triggered only by chase choice")}
 for _,id:=range []int{101,102,104} {
  when,err:=time.Parse("15:04",byID[id])
  if err!=nil || !when.Before(chase) {t.Errorf("event %d must precede rooftop chase",id)}
 }
}

package cases

import (
 "path/filepath"
 "testing"
 "time"

 "detective/internal/game"
)

func TestExpress1928RooftopChaseBranchesAndReturnsSafely(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 makeState:=func()*game.GameState {
  return &game.GameState{
   Events:def.Events,Characters:def.Characters,
   FoundFacts:map[int]bool{},TriggeredEvents:map[int]bool{},
  }
 }
 t.Run("pursue and identify Graf",func(t *testing.T){
  s:=makeState()
  if !game.TriggerEvent(s,103) {t.Fatal("chase opening unavailable")}
  if !game.ChooseEventOption(s,1) {t.Fatal("cannot follow hooded figure")}
  if s.ActiveEvent==nil || *s.ActiveEvent!=105 {t.Fatal("rooftop follow-up not triggered")}
  if !s.TriggeredEvents[105] {t.Fatal("follow-up event not recorded")}
  if !game.ChooseEventOption(s,1) {t.Fatal("cannot retreat with letter")}
  if s.ActiveEvent!=nil || !s.FoundFacts[18] {t.Fatal("chase resolution lost optional evidence")}
 })
 t.Run("refuse rooftop pursuit",func(t *testing.T){
  s:=makeState()
  if !game.TriggerEvent(s,103) || !game.ChooseEventOption(s,4) {t.Fatal("safe refusal unavailable")}
  if s.ActiveEvent!=nil || s.TriggeredEvents[105] {t.Fatal("rooftop chase incorrectly forced")}
 })
 t.Run("tunnel deadline forces retreat",func(t *testing.T){
  s:=makeState()
  if !game.TriggerEvent(s,103) || !game.ChooseEventOption(s,1) {t.Fatal("cannot enter rooftop chase")}
  s.EventDeadline=time.Now().Add(-time.Second)
  if !game.ExpireActiveEvent(s) {t.Fatal("expired rooftop decision did not resolve")}
  if s.ActiveEvent!=nil {t.Fatal("rooftop event left blocking state")}
  if s.FoundFacts[18] {t.Fatal("timeout awarded letter evidence")}
 })
}

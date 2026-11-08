package game

import (
 "testing"
 "time"
)

func TestSceneChoiceMovesAndChargesTime(t *testing.T) {
 start:=time.Date(2026,1,1,22,0,0,0,time.UTC)
 state:=&GameState{
  Clock:NewClock(start),CurrentScene:"corridor",
  Scenes:[]Scene{
   {ID:"corridor",Choices:[]Choice{{ID:"go",Kind:"action",NextScene:"compartment",Action:Action{TimeCost:30}}}},
   {ID:"compartment"},
  },
 }
 if err:=ChooseSceneOption(state,"go");err!=nil { t.Fatal(err) }
 if state.CurrentScene!="compartment" { t.Fatalf("scene=%s",state.CurrentScene) }
 if got:=state.Clock.GetTime().Sub(start);got!=30*time.Second { t.Fatalf("elapsed=%v",got) }
}

func TestSceneChoiceCannotBypassRequirements(t *testing.T) {
 start:=time.Date(2026,1,1,22,0,0,0,time.UTC)
 state:=&GameState{
  Clock:NewClock(start),CurrentScene:"corridor",
  Scenes:[]Scene{{ID:"corridor",Choices:[]Choice{{ID:"locked",Action:Action{
   Requirements:[]Requirement{{Type:"fact_found",ID:999}},TimeCost:30,
  }}}}},
 }
 if err:=ChooseSceneOption(state,"locked");err==nil { t.Fatal("expected unavailable choice") }
 if !state.Clock.GetTime().Equal(start) { t.Fatal("failed choice changed time") }
 view,err:=CurrentSceneView(state)
 if err!=nil { t.Fatal(err) }
 if len(view.Choices)!=0 { t.Fatal("locked choice exposed to client") }
}

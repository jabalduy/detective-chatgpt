package server

import (
 "detective/internal/game"
 "testing"
 "time"
)

func TestRestoreActiveConversationAndProgress(t *testing.T){
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 original:=&game.GameState{
  Case:game.CaseInfo{ID:"case_train"},
  Clock:game.NewClock(start),CurrentScene:"carriage",
  ActiveConversation:"witness",ConversationNode:"question",
  Characters:map[int]game.CharacterState{7:{Location:"carriage",Present:true,Trust:3}},
  FoundClues:map[int]bool{11:true},FoundFacts:map[int]bool{22:true},
  FoundMotives:map[int]bool{33:true},
  OpenLocations:map[int]bool{2:true},OpenObjects:map[int]bool{4:true},
  SearchedObjects:map[int]bool{4:true},
  ProcessedStations:map[string]bool{"stop:arrival":true},
  TriggeredEvents:map[int]bool{9:true},
  ExaminationStarted:map[string]time.Time{"lab":start},
  ExaminationCompleted:map[string]bool{"lab":true},
  PresentedEvidence:map[string]bool{"presentation":true},
 }
 original.Clock.Pause()
 saved:=snapshot(original)
 original.Clock.StopTime()
 restored:=&game.GameState{Clock:game.NewClock(start)}
 applySave(restored,saved)
 defer restored.Clock.StopTime()
 if restored.ActiveConversation!="witness" || restored.ConversationNode!="question" || !restored.Clock.IsPaused(){
  t.Fatal("conversation or paused clock lost")
 }
 if restored.Characters[7].Trust!=3 || !restored.FoundClues[11] || !restored.FoundFacts[22] || !restored.FoundMotives[33]{
  t.Fatal("character trust or evidence lost")
 }
 if !restored.OpenLocations[2] || !restored.OpenObjects[4] || !restored.SearchedObjects[4]{
  t.Fatal("exploration progress lost")
 }
 if !restored.ProcessedStations["stop:arrival"] || !restored.TriggeredEvents[9] || !restored.ExaminationCompleted["lab"] || !restored.PresentedEvidence["presentation"]{
  t.Fatal("event or examination progress lost")
 }
}

func TestRestoreActiveTimedEventAndPausedClock(t *testing.T){
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 id:=8
 original:=&game.GameState{
  Case:game.CaseInfo{ID:"case_train"},Clock:game.NewClock(start),
  ActiveEvent:&id,EventDeadline:time.Now().Add(20*time.Second),
  TriggeredEvents:map[int]bool{8:true},
 }
 original.Clock.Pause()
 saved:=snapshot(original)
 original.Clock.StopTime()
 restored:=&game.GameState{Clock:game.NewClock(start)}
 applySave(restored,saved)
 defer restored.Clock.StopTime()
 if restored.ActiveEvent==nil || *restored.ActiveEvent!=8 || !restored.TriggeredEvents[8]{
  t.Fatal("active event or triggered flag lost")
 }
 if !restored.Clock.IsPaused(){t.Fatal("paused clock resumed during restore")}
 if remaining:=time.Until(restored.EventDeadline);remaining<18*time.Second || remaining>21*time.Second{
  t.Fatalf("event remaining time changed: %v",remaining)
 }
}

func TestRestoreNilMapsAllowsNewEvidence(t *testing.T){
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 original:=&game.GameState{Case:game.CaseInfo{ID:"case_train"},Clock:game.NewClock(start)}
 saved:=snapshot(original)
 original.Clock.StopTime()
 restored:=&game.GameState{Clock:game.NewClock(start)}
 applySave(restored,saved)
 defer restored.Clock.StopTime()
 if !game.PerformAction(restored,game.Action{Effects:[]game.Effect{{Type:"discover_fact",ID:55}}}){
  t.Fatal("restored game refused action")
 }
 if !restored.FoundFacts[55]{t.Fatal("new fact not recorded")}
}

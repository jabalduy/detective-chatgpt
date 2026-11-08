package game

import (
 "testing"
 "time"
)

func TestFinishedGameDoesNotAdvanceStationsOrExaminations(t *testing.T){
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 clock:=NewClock(start)
 defer clock.StopTime()
 clock.Add(20*time.Minute)
 s:=&GameState{
  Clock:clock,RouteStart:start,
  Ending:&GameEnding{Type:"win"},
  Stations:[]Station{{ID:"station",Arrival:"23:55",ArriveEffects:[]Effect{{Type:"discover_fact",ID:10}}}},
  Examinations:[]Examination{{ID:"lab",DurationSeconds:60,Effects:[]Effect{{Type:"discover_fact",ID:20}}}},
  ExaminationStarted:map[string]time.Time{"lab":start},
  FoundFacts:map[int]bool{},
 }
 TickTrain(s)
 TickExaminations(s)
 ExaminationViews(s)
 if s.FoundFacts[10] || s.FoundFacts[20] || s.ProcessedStations["station:arrival"] || s.ExaminationCompleted["lab"]{
  t.Fatal("finished game advanced scheduled effects")
 }
}

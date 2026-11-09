package game

import (
 "testing"
 "time"
)

func TestHintEscalationAndUsage(t *testing.T){
 want:=[]int{180,300,480,720,900,900}
 for i,cost:=range want {if got:=HintCostSeconds(i);got!=cost{t.Fatalf("hint %d cost %d want %d",i,got,cost)}}
 start:=time.Date(2026,1,1,23,0,0,0,time.UTC)
 s:=&GameState{Clock:NewClock(start),Hints:[]Hint{{ID:"first",Text:"Look at the timetable"},{ID:"second",Text:"Check the ticket"}}}
 defer s.Clock.StopTime()
 first,err:=TakeHint(s)
 if err!=nil || first.ID!="first" || first.CostSeconds!=180 {t.Fatalf("first hint: %+v %v",first,err)}
 second,err:=TakeHint(s)
 if err!=nil || second.ID!="second" || second.CostSeconds!=300 {t.Fatalf("second hint: %+v %v",second,err)}
 if got:=s.Clock.GetTime();!got.Equal(start.Add(480*time.Second)){t.Fatalf("clock time %s",got)}
 if _,err:=TakeHint(s);err==nil {t.Fatal("reused exhausted hint")}
 s.InvestigationDeadline=time.Now().Add(-time.Second)
 if _,err:=AvailableHint(s);err==nil {t.Fatal("hint available after investigation")}
}

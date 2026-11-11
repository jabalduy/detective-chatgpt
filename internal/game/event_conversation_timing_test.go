package game

import (
 "testing"
 "time"
)

func TestTimedEventWaitsForConversationAndResumesAfterward(t *testing.T) {
 start:=time.Date(1928,time.October,1,21,0,0,0,time.UTC)
 clock:=NewClock(start)
 defer clock.StopTime()
 s:=&GameState{
  Clock:clock,RouteStart:start,CurrentScene:"salon",
  Characters:map[int]CharacterState{2:{Location:"salon",Present:true}},
  Conversations:[]Conversation{{ID:"boris",CharacterID:2,StartNode:"hello",Nodes:[]ConversationNode{{
   ID:"hello",Choices:[]ConversationChoice{{ID:"leave",End:true}},
  }}}},
  Events:[]Event{{ID:101,At:"21:12",TimeLimit:10,DefaultOptionID:2,Options:[]EventOption{
   {ID:1,Action:Action{Effects:[]Effect{{Type:"discover_fact",ID:17}}}},
   {ID:2},
  }}},
  FoundFacts:map[int]bool{},TriggeredEvents:map[int]bool{},
 }
 clock.Add(11*time.Minute+59*time.Second)
 if err:=StartConversation(s,"boris");err!=nil {t.Fatal(err)}
 if !clock.IsPaused() {t.Fatal("route clock should pause for dialogue")}
 clock.Add(2*time.Second) // Simulate an action crossing the event threshold.
 TickEvents(s)
 if s.ActiveEvent!=nil {t.Fatal("event interrupted active conversation")}
 if s.TriggeredEvents[101] {t.Fatal("deferred event marked as triggered")}
 if _,err:=ChooseConversationOption(s,"leave");err!=nil {t.Fatal(err)}
 if clock.IsPaused() {t.Fatal("route clock did not resume after conversation")}
 TickEvents(s)
 if s.ActiveEvent==nil || *s.ActiveEvent!=101 {t.Fatal("overdue event not triggered after conversation")}
 if !s.TriggeredEvents[101] {t.Fatal("triggered event not recorded")}
 if err:=StartConversation(s,"boris");err==nil {t.Fatal("conversation started during active event")}
 if err:=ChooseSceneOption(s,"anything");err==nil {t.Fatal("scene interaction started during active event")}
 if !ChooseEventOption(s,2) {t.Fatal("cannot resolve deferred event")}
 if s.ActiveEvent!=nil {t.Fatal("event still active after resolution")}
 TickEvents(s)
 if s.ActiveEvent!=nil {t.Fatal("event retriggered after resolution")}
}

func TestTimedEventExpiryUsesDefaultAndUnblocksInvestigation(t *testing.T) {
 start:=time.Date(1928,time.October,1,21,0,0,0,time.UTC)
 clock:=NewClock(start)
 defer clock.StopTime()
 s:=&GameState{
  Clock:clock,RouteStart:start,CurrentScene:"corridor",
  Events:[]Event{{ID:101,At:"21:12",TimeLimit:10,DefaultOptionID:2,Options:[]EventOption{
   {ID:1,Action:Action{Effects:[]Effect{{Type:"discover_fact",ID:17}}}},
   {ID:2,Action:Action{Effects:[]Effect{{Type:"discover_fact",ID:19}}}},
  }}},
  FoundFacts:map[int]bool{},TriggeredEvents:map[int]bool{},
 }
 clock.Add(12*time.Minute)
 TickEvents(s)
 if s.ActiveEvent==nil {t.Fatal("timed event did not start")}
 s.EventDeadline=time.Now().Add(-time.Second)
 TickEvents(s)
 if s.ActiveEvent!=nil {t.Fatal("expired event not cleared")}
 if !s.FoundFacts[19] {t.Fatal("default action not applied")}
 if s.FoundFacts[17] {t.Fatal("unselected option applied")}
 TickEvents(s)
 if s.ActiveEvent!=nil {t.Fatal("expired event retriggered")}
}

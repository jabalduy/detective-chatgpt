package game

import "time"

type Event struct {
 ID int `json:"id"`
 Title string `json:"title"`
 Description string `json:"description"`
 TimeLimit int `json:"time_limit"`
 At string `json:"at,omitempty"`
 DefaultOptionID int `json:"default_option_id,omitempty"`
 Options []EventOption `json:"options"`
}
type EventOption struct {
 ID int `json:"id"`
 Text string `json:"text"`
 Action
}

func TriggerEvent(state *GameState,eventID int) bool {
 if state.Ending!=nil {return false}
 if state.ActiveEvent!=nil || state.ActiveConversation!="" || state.ActiveDialogue!=nil {return false}
 for _,event:=range state.Events {
  if event.ID!=eventID {continue}
  state.ActiveEvent=&eventID
  if event.TimeLimit>0 {state.EventDeadline=time.Now().Add(time.Duration(event.TimeLimit)*time.Second)} else {state.EventDeadline=time.Time{}}
  if state.TriggeredEvents==nil {state.TriggeredEvents=make(map[int]bool)}
  state.TriggeredEvents[eventID]=true
  return true
 }
 return false
}
func eventByID(s *GameState,id int)*Event{
 for i:=range s.Events {if s.Events[i].ID==id{return &s.Events[i]}}
 return nil
}
func ChooseEventOption(state *GameState,optionID int)bool{
 if state.ActiveEvent==nil{return false}
 event:=eventByID(state,*state.ActiveEvent)
 if event==nil{return false}
 if !state.EventDeadline.IsZero() && !time.Now().Before(state.EventDeadline) {ExpireActiveEvent(state);return false}
 for _,option:=range event.Options {
  if option.ID!=optionID {continue}
  if !PerformAction(state,option.Action){return false}
  state.ActiveEvent=nil
  state.EventDeadline=time.Time{}
  return true
 }
 return false
}
// Expiration is enforced on the server; the client cannot dismiss a timed event.
func ExpireActiveEvent(s *GameState)bool{
 if s.ActiveEvent==nil{return false}
 e:=eventByID(s,*s.ActiveEvent)
 if e==nil{return false}
 if e.TimeLimit<=0 || time.Now().Before(s.EventDeadline){return false}
 // Clear first so effects can trigger another event.
 s.ActiveEvent=nil
 s.EventDeadline=time.Time{}
 for _,option:=range e.Options {
  if option.ID==e.DefaultOptionID && e.DefaultOptionID!=0 {
   PerformAction(s,option.Action)
   break
  }
 }
 return true
}
// TickEvents checks story-clock triggers; pending events wait until dialogue ends.
func TickEvents(s *GameState){
 if s.Ending!=nil{return}
 ExpireActiveEvent(s)
 if s.ActiveEvent!=nil || s.ActiveConversation!="" || s.ActiveDialogue!=nil || s.Clock==nil {return}
 now:=s.Clock.GetTime()
 for _,e:=range s.Events {
  if e.At=="" || s.TriggeredEvents[e.ID] {continue}
  when,err:=routeTime(s.RouteStart,e.At)
  if err==nil && !now.Before(when) {TriggerEvent(s,e.ID);return}
 }
}

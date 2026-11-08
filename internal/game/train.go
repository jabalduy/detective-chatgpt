package game

import "time"

// Station times use 24h HH:MM, including overnight routes.
type Station struct {
 ID string `json:"id"`
 Name string `json:"name"`
 Arrival string `json:"arrival"`
 Departure string `json:"departure"`
 ArriveEffects []Effect `json:"arrive_effects,omitempty"`
 DepartEffects []Effect `json:"depart_effects,omitempty"`
}
type TrainStatus struct {
 CurrentTime string `json:"current_time"`
 AtStation string `json:"at_station,omitempty"`
 NextStation string `json:"next_station,omitempty"`
 NextArrival string `json:"next_arrival,omitempty"`
}

// TickTrain must be called with the session locked. It catches up after time jumps.
func TickTrain(s *GameState) {
 if s.Clock==nil || s.Ending!=nil {return}
 now:=s.Clock.GetTime()
 if s.ProcessedStations==nil {s.ProcessedStations=make(map[string]bool)}
 for _,st:=range s.Stations {
  for _,phase:=range []struct{suffix,clock string;effects []Effect}{
   {"arrival",st.Arrival,st.ArriveEffects},{"departure",st.Departure,st.DepartEffects},
  } {
   if phase.clock=="" {continue}
   when,err:=routeTime(s.RouteStart,phase.clock)
   if err!=nil || now.Before(when) {continue}
   key:=st.ID+":"+phase.suffix
   if s.ProcessedStations[key] {continue}
   s.ProcessedStations[key]=true
   ApplyEffects(s,phase.effects)
  }
 }
}

func TrainView(s *GameState) TrainStatus {
 v:=TrainStatus{}
 if s.Clock==nil {return v}
 now:=s.Clock.GetTime()
 v.CurrentTime=now.Format("15:04:05")
 for _,st:=range s.Stations {
  arrival,errA:=routeTime(s.RouteStart,st.Arrival)
  departure,errD:=routeTime(s.RouteStart,st.Departure)
  if errA!=nil || errD!=nil {continue}
  if !now.Before(arrival)&&now.Before(departure) {v.AtStation=st.Name}
  if now.Before(arrival) && (v.NextStation=="" || arrival.Before(nextArrival(s,v.NextArrival))) {
   v.NextStation=st.Name
   v.NextArrival=arrival.Format(time.RFC3339)
  }
 }
 return v
}
func nextArrival(s *GameState,raw string)time.Time {
 t,_:=time.Parse(time.RFC3339,raw);return t
}
func routeTime(start time.Time,clock string)(time.Time,error){
 parsed,err:=time.Parse("15:04",clock)
 if err!=nil{return time.Time{},err}
 when:=time.Date(start.Year(),start.Month(),start.Day(),parsed.Hour(),parsed.Minute(),0,0,start.Location())
 if when.Before(start){when=when.AddDate(0,0,1)}
 return when,nil
}

package server

import (
 "detective/internal/game"
 "encoding/json"
 "net/http"
 "time"
)

// gameStateView is a read-only snapshot for a single-operator frontend.
// It intentionally excludes hidden evidence, action effects, solution data,
// and unrevealed event definitions.
type gameStateView struct {
 CaseID string `json:"case_id"`
 Started bool `json:"started"`
 Phase game.RealTimePhase `json:"phase"`
 NarrativeAct game.NarrativeAct `json:"narrative_act"`
 InvestigationDeadline time.Time `json:"investigation_deadline,omitempty"`
 AccusationDeadline time.Time `json:"accusation_deadline,omitempty"`
 Time string `json:"time,omitempty"`
 Paused bool `json:"paused"`
 Scene *game.SceneView `json:"scene,omitempty"`
 Train *game.TrainStatus `json:"train,omitempty"`
 Characters map[int]game.CharacterState `json:"characters,omitempty"`
 Conversation *game.ConversationView `json:"conversation,omitempty"`
 ActiveEvent *activeEventView `json:"active_event,omitempty"`
 Ending *game.GameEnding `json:"ending,omitempty"`
 FoundClueIDs []int `json:"found_clue_ids"`
 FoundFactIDs []int `json:"found_fact_ids"`
 FoundMotiveIDs []int `json:"found_motive_ids"`
}
type activeEventView struct {
 ID int `json:"id"`
 Title string `json:"title"`
 Description string `json:"description"`
 Options []eventOptionView `json:"options"`
 RemainingSeconds int `json:"remaining_seconds,omitempty"`
}
type eventOptionView struct {
 ID int `json:"id"`
 Text string `json:"text"`
}
func sortedFound(m map[int]bool)[]int{
 ids:=make([]int,0)
 for id,found:=range m{if found{ids=append(ids,id)}}
 // Stable ordering makes snapshots easier to compare and render.
 for i:=1;i<len(ids);i++{for j:=i;j>0 && ids[j]<ids[j-1];j--{ids[j],ids[j-1]=ids[j-1],ids[j]}}
 return ids
}
func buildGameStateView(s *game.GameState)(gameStateView,error){
 v:=gameStateView{
  CaseID:s.Case.ID,Started:true,
  Phase:s.PhaseAt(time.Now()),NarrativeAct:s.CurrentNarrativeAct(),InvestigationDeadline:s.InvestigationDeadline,AccusationDeadline:s.AccusationDeadline,
  FoundClueIDs:sortedFound(s.FoundClues),
  FoundFactIDs:sortedFound(s.FoundFacts),
  FoundMotiveIDs:sortedFound(s.FoundMotives),
  Characters:make(map[int]game.CharacterState,len(s.Characters)),
  Ending:s.Ending,
 }
 for id,c:=range s.Characters{v.Characters[id]=c}
 if s.Clock!=nil {v.Time=s.Clock.GetTime().Format(time.RFC3339);v.Paused=s.Clock.IsPaused();train:=game.TrainView(s);v.Train=&train}
 if len(s.Scenes)>0 {
  scene,err:=game.CurrentSceneView(s);if err!=nil{return v,err}
  v.Scene=&scene
 }
 if s.ActiveConversation!=""{
  conversation,err:=game.CurrentConversationView(s);if err!=nil{return v,err}
  v.Conversation=&conversation
 }
 if s.ActiveEvent!=nil {
  for _,e:=range s.Events {if e.ID==*s.ActiveEvent{
   ev:=&activeEventView{ID:e.ID,Title:e.Title,Description:e.Description,Options:make([]eventOptionView,0,len(e.Options))}
   for _,option:=range e.Options{if game.CanPerformAction(s,option.Action){ev.Options=append(ev.Options,eventOptionView{ID:option.ID,Text:option.Text})}}
   if !s.EventDeadline.IsZero(){remaining:=time.Until(s.EventDeadline);if remaining>0{ev.RemainingSeconds=int((remaining+time.Second-1)/time.Second)}}
   v.ActiveEvent=ev
   break
  }}
 }
 return v,nil
}

func stateHandler(w http.ResponseWriter,r *http.Request){
 if r.Method!=http.MethodGet{http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return}
 session,err:=getStartedGame(w,r)
 if err!=nil{http.Error(w,"game not started",http.StatusConflict);return}
 session.Mu.Lock()
 defer session.Mu.Unlock()
 game.TickTrain(session.State)
 game.TickEvents(session.State)
 view,err:=buildGameStateView(session.State)
 if err!=nil{http.Error(w,err.Error(),http.StatusInternalServerError);return}
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(view)
}

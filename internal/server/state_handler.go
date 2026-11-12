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
type availableConversationView struct { ID string `json:"id"`; CharacterID int `json:"character_id"`; Name string `json:"name"` }
type gameStateView struct {
 CurrentScene string `json:"current_scene,omitempty"`
 AvailableConversations []availableConversationView `json:"available_conversations"`
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
  CaseID:s.Case.ID,Started:true,CurrentScene:s.CurrentScene,
  AvailableConversations:make([]availableConversationView,0),
  Phase:s.PhaseAt(time.Now()),NarrativeAct:s.CurrentNarrativeAct(),InvestigationDeadline:s.InvestigationDeadline,AccusationDeadline:s.AccusationDeadline,
  FoundClueIDs:sortedFound(s.FoundClues),
  FoundFactIDs:sortedFound(s.FoundFacts),
  FoundMotiveIDs:sortedFound(s.FoundMotives),
  Characters:make(map[int]game.CharacterState,len(s.Characters)),
  Ending:s.Ending,
 }
 // The prologue must not disclose investigation evidence IDs, including
 // IDs restored from an inconsistent or older save.
 if s.CurrentNarrativeAct()==game.ActPrologue {
  v.FoundClueIDs=[]int{}
  v.FoundFactIDs=[]int{}
  v.FoundMotiveIDs=[]int{}
 }
 for id,c:=range s.Characters{v.Characters[id]=c}
 // Only publish conversations that can be started here and in this act.
 if s.ActiveEvent==nil && s.ActiveConversation=="" && s.ActiveDialogue==nil {
  for _,conv:=range s.Conversations {
   if !game.ConversationAvailableInAct(s,&conv) {continue}
   ch,ok:=s.Characters[conv.CharacterID]
   if !ok || !ch.Present || ch.Location!=s.CurrentScene {continue}
   names:=map[int]string{1:"Марта",2:"Борис",3:"Софи",4:"Граф Антон",5:"Клара",6:"Макс",7:"Роза",8:"Лео"}
   name:=names[conv.CharacterID]
   if name=="" {name="Пассажир"}
   v.AvailableConversations=append(v.AvailableConversations,availableConversationView{ID:conv.ID,CharacterID:conv.CharacterID,Name:name})
  }
 }
 if s.Clock!=nil {v.Time=s.Clock.GetTime().Format(time.RFC3339);v.Paused=s.Clock.IsPaused();train:=game.TrainView(s);v.Train=&train}
 if s.CurrentNarrativeAct()==game.ActInvestigation && len(s.Scenes)>0 {
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
 // Polling /state advances the narrative once the official discovery is due.
 // BeginInvestigation keeps urgent interactions intact and is idempotent.
 // GET /state is already covered by saveAfterMutation.
 if session.State.CurrentNarrativeAct()==game.ActPrologue {
  _,transitionErr:=session.State.BeginInvestigation(time.Now())
  if transitionErr!=nil && transitionErr!=game.ErrDiscoveryTooEarly && transitionErr!=game.ErrDiscoveryInteractionActive {
   http.Error(w,transitionErr.Error(),http.StatusInternalServerError);return
  }
 }
 view,err:=buildGameStateView(session.State)
 if err!=nil{http.Error(w,err.Error(),http.StatusInternalServerError);return}
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(view)
}

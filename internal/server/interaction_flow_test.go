package server

import (
 "detective/internal/game"
 "encoding/json"
 "net/http"
 "testing"
 "time"
)

func TestHTTPSceneConversationEventFlow(t *testing.T){
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"test"},
  Clock:game.NewClock(start),RouteStart:start,CurrentScene:"corridor",
  Scenes:[]game.Scene{
   {ID:"corridor",Choices:[]game.Choice{{ID:"enter",NextScene:"carriage"}}},
   {ID:"carriage",Choices:[]game.Choice{{ID:"leave",NextScene:"corridor"}}},
  },
  Characters:map[int]game.CharacterState{7:{Location:"carriage",Present:true}},
  Conversations:[]game.Conversation{{ID:"witness",CharacterID:7,StartNode:"first",Nodes:[]game.ConversationNode{{
   ID:"first",Choices:[]game.ConversationChoice{{ID:"finish",End:true,Action:game.Action{Effects:[]game.Effect{{Type:"discover_clue",ID:11}}}}},
  }}}},
  Events:[]game.Event{{ID:3,TimeLimit:60,Options:[]game.EventOption{{ID:4,Text:"Respond",Action:game.Action{Effects:[]game.Effect{{Type:"discover_fact",ID:12}}}}}}},
  FoundClues:map[int]bool{},FoundFacts:map[int]bool{},
 }
 defer session.State.Clock.StopTime()
 check:=func(method,path string,want int)*httptestResponse{
  t.Helper()
  r:=requestWithCookie(handler,cookie,method,path,nil)
  if r.Code!=want{t.Fatalf("%s %s: got %d want %d; %s",method,path,r.Code,want,r.Body.String())}
  return &httptestResponse{body:r.Body.Bytes()}
 }
 check(http.MethodPost,"/scene/choose?choiceID=enter",200)
 check(http.MethodPost,"/conversation/start?id=witness",200)
 if !session.State.Clock.IsPaused(){t.Fatal("conversation did not pause clock")}
 check(http.MethodPost,"/scene/choose?choiceID=leave",409)
 check(http.MethodPost,"/conversation/choose?choiceID=finish",200)
 if session.State.Clock.IsPaused(){t.Fatal("clock still paused after conversation")}
 if !session.State.FoundClues[11]{t.Fatal("conversation clue missing")}
 if !game.TriggerEvent(session.State,3){t.Fatal("event not triggered")}
 check(http.MethodPost,"/scene/choose?choiceID=leave",409)
 r:=check(http.MethodGet,"/state",200)
 var state gameStateView
 if err:=json.Unmarshal(r.body,&state);err!=nil{t.Fatal(err)}
 if state.ActiveEvent==nil || state.ActiveEvent.ID!=3{t.Fatal("active event missing from state")}
 check(http.MethodPost,"/event/choose?optionID=4",200)
 if session.State.ActiveEvent!=nil || !session.State.FoundFacts[12]{t.Fatal("event resolution failed")}
 check(http.MethodPost,"/scene/choose?choiceID=leave",200)
 if session.State.CurrentScene!="corridor"{t.Fatal("scene did not change")}
 saved,err:=readSave(savePath(cookie.Value))
 if err!=nil{t.Fatal(err)}
 if saved.CurrentScene!="corridor" || !saved.FoundClues[11] || !saved.FoundFacts[12]{
  t.Fatal("completed flow not autosaved")
 }
}

type httptestResponse struct{body []byte}

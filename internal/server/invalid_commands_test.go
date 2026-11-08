package server

import (
 "detective/internal/game"
 "net/http"
 "testing"
 "time"
)

func TestInvalidSceneCommandsLeaveStateUnchanged(t *testing.T){
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"test"},Clock:game.NewClock(start),
  CurrentScene:"corridor",
  Scenes:[]game.Scene{
   {ID:"corridor",Choices:[]game.Choice{
    {ID:"locked",NextScene:"carriage",Action:game.Action{Requirements:[]game.Requirement{{Type:"clue_found",ID:99}},Effects:[]game.Effect{{Type:"discover_fact",ID:12}}}},
    {ID:"broken",NextScene:"missing",Action:game.Action{Effects:[]game.Effect{{Type:"discover_fact",ID:13}}}},
   }},
   {ID:"carriage"},
  },
  FoundClues:map[int]bool{},FoundFacts:map[int]bool{},
 }
 defer session.State.Clock.StopTime()
 for _,tc:=range []struct{method,path string;want int}{
  {http.MethodGet,"/scene/choose?choiceID=locked",405},
  {http.MethodPost,"/scene/choose",400},
  {http.MethodPost,"/scene/choose?choiceID=unknown",409},
  {http.MethodPost,"/scene/choose?choiceID=locked",409},
  {http.MethodPost,"/scene/choose?choiceID=broken",409},
 }{
  response:=requestWithCookie(handler,cookie,tc.method,tc.path,nil)
  if response.Code!=tc.want{t.Fatalf("%s %s: got %d want %d: %s",tc.method,tc.path,response.Code,tc.want,response.Body.String())}
  if session.State.CurrentScene!="corridor" || session.State.FoundFacts[12] || session.State.FoundFacts[13]{
   t.Fatalf("rejected command changed game state: %s",tc.path)
  }
 }
}

func TestInvalidConversationCommandsLeaveStateUnchanged(t *testing.T){
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"test"},Clock:game.NewClock(start),
  CurrentScene:"carriage",
  Characters:map[int]game.CharacterState{7:{Location:"carriage",Present:true}},
  Conversations:[]game.Conversation{{ID:"witness",CharacterID:7,StartNode:"first",Nodes:[]game.ConversationNode{{
   ID:"first",Choices:[]game.ConversationChoice{
    {ID:"secret",End:true,Action:game.Action{Requirements:[]game.Requirement{{Type:"clue_found",ID:99}},Effects:[]game.Effect{{Type:"discover_fact",ID:12}}}},
    {ID:"broken",NextNode:"missing",Action:game.Action{Effects:[]game.Effect{{Type:"discover_fact",ID:13}}}},
   },
  }}}},
  FoundClues:map[int]bool{},FoundFacts:map[int]bool{},
 }
 defer session.State.Clock.StopTime()
 r:=requestWithCookie(handler,cookie,http.MethodPost,"/conversation/start?id=witness",nil)
 if r.Code!=200{t.Fatalf("start: %d %s",r.Code,r.Body.String())}
 for _,id:=range []string{"unknown","secret","broken"}{
  r=requestWithCookie(handler,cookie,http.MethodPost,"/conversation/choose?choiceID="+id,nil)
  if r.Code!=409{t.Fatalf("choice %s: got %d want 409: %s",id,r.Code,r.Body.String())}
  if session.State.ActiveConversation!="witness" || session.State.ConversationNode!="first" || !session.State.Clock.IsPaused() || session.State.FoundFacts[12] || session.State.FoundFacts[13]{
   t.Fatalf("rejected conversation choice changed game state: %s",id)
  }
 }
}

package server

import (
 "bytes"
 "detective/internal/game"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "testing"
 "time"
)

func testSession(t *testing.T)(http.Handler,*http.Cookie,*Session){
 t.Helper()
 handler:=NewHandler()
 req:=httptest.NewRequest(http.MethodGet,"/scene",nil)
 recorder:=httptest.NewRecorder()
 session,err:=getGame(recorder,req)
 if err!=nil{t.Fatal(err)}
 cookies:=recorder.Result().Cookies()
 if len(cookies)==0{t.Fatal("session cookie not created")}
 t.Cleanup(func(){gamesMu.Lock();delete(games,cookies[0].Value);gamesMu.Unlock()})
 return handler,cookies[0],session
}
func requestWithCookie(handler http.Handler,cookie *http.Cookie,method,path string,body []byte)*httptest.ResponseRecorder{
 req:=httptest.NewRequest(method,path,bytes.NewReader(body))
 req.AddCookie(cookie)
 recorder:=httptest.NewRecorder()
 handler.ServeHTTP(recorder,req)
 return recorder
}
func TestHTTPGameFlow(t *testing.T){
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 session.State=&game.GameState{
  Clock:game.NewClock(start),RouteStart:start,CurrentScene:"corridor",
  Scenes:[]game.Scene{{ID:"corridor",Choices:[]game.Choice{{ID:"inspect",Action:game.Action{Effects:[]game.Effect{{Type:"discover_clue",ID:1}}}}}}},
  FoundClues:make(map[int]bool),FoundFacts:make(map[int]bool),FoundMotives:make(map[int]bool),
  OpenDialogues:make(map[game.DialogueID]bool),
  Characters:map[int]game.CharacterState{1:{Location:"corridor",Present:true}},
  Conversations:[]game.Conversation{{ID:"witness",CharacterID:1,StartNode:"first",Nodes:[]game.ConversationNode{
   {ID:"first",Text:"Hello",Choices:[]game.ConversationChoice{{ID:"finish",End:true,Action:game.Action{Effects:[]game.Effect{{Type:"discover_motive",ID:2}}}}}},
  }}},
  Examinations:[]game.Examination{{ID:"lab",DurationSeconds:10,Requirements:[]game.Requirement{{Type:"clue_found",ID:1}},Result:"Matched",Effects:[]game.Effect{{Type:"discover_fact",ID:3}}}},
  Motives:[]game.Motive{{ID:2,SuspectID:1}},
  Case:game.CaseInfo{ID:"test-case",Solution:game.Solution{SuspectID:1,MotiveID:2,FactIDs:[]int{3}},Endings:game.Endings{Win:"Solved"}},
 }
 check:=func(method,path string,body []byte,want int)*httptest.ResponseRecorder{
  t.Helper()
  r:=requestWithCookie(handler,cookie,method,path,body)
  if r.Code!=want{t.Fatalf("%s %s: status %d, want %d; body=%s",method,path,r.Code,want,r.Body.String())}
  return r
 }
 check("GET","/scene",nil,200)
 check("POST","/scene/choose?choiceID=inspect",nil,200)
 check("POST","/examinations/start?id=lab",nil,200)
 check("POST","/conversation/start?id=witness",nil,200)
 check("POST","/scene/choose?choiceID=inspect",nil,409)
 check("POST","/conversation/choose?choiceID=finish",nil,200)
 session.State.Clock.Add(10*time.Second)
 response:=check("GET","/examinations",nil,200)
 var exams []game.ExaminationProgress
 if err:=json.Unmarshal(response.Body.Bytes(),&exams);err!=nil||len(exams)!=1||exams[0].Status!="completed"{t.Fatalf("examinations=%+v err=%v",exams,err)}
 payload:=[]byte(`{"suspect_id":1,"motive_id":2,"fact_ids":[3]}`)
 result:=check("POST","/accuse",payload,200)
 var verdict struct{TypeEnd string `json:"typeEnd"`}
 if err:=json.Unmarshal(result.Body.Bytes(),&verdict);err!=nil||verdict.TypeEnd!="win"{t.Fatalf("verdict=%+v err=%v",verdict,err)}
 check("GET","/ending",nil,200)
 check("POST","/scene/choose?choiceID=inspect",nil,409)
 check("POST","/accuse",payload,409)
 if !session.State.Clock.IsPaused(){t.Fatal("clock must stop after ending")}
 saved,err:=readSave(savePath(cookie.Value))
 if err!=nil{t.Fatalf("final autosave: %v",err)}
 if saved.Ending==nil || saved.Ending.Type!="win" || !saved.ClockPaused{
  t.Fatalf("final verdict or paused clock not saved: %+v",saved)
 }
 if !saved.FoundClues[1] || !saved.FoundFacts[3] || !saved.FoundMotives[2] || !saved.ExaminationCompleted["lab"]{
  t.Fatal("full investigation progress missing from final save")
 }
 restored:=&game.GameState{Clock:game.NewClock(start)}
 applySave(restored,saved)
 defer restored.Clock.StopTime()
 if restored.Ending==nil || restored.Ending.Type!="win" || !restored.Clock.IsPaused(){
  t.Fatal("final result not restored correctly")
 }
 if _,err:=game.SubmitAccusation(restored,game.Accusation{SuspectID:1,MotiveID:2,FactIDs:[]int{3}});err==nil{
  t.Fatal("restored finished game accepted another accusation")
 }
}
func TestHTTPRequiresStartedGame(t *testing.T){
 handler,cookie,_:=testSession(t)
 for _,path:=range []string{"/scene","/conversation","/examinations","/ending"}{
  response:=requestWithCookie(handler,cookie,"GET",path,nil)
  if response.Code!=http.StatusConflict{t.Fatalf("%s returned %d",path,response.Code)}
 }
}

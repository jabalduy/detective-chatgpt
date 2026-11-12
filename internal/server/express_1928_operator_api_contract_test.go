package server

import (
 "encoding/json"
 "net/http"
 "path/filepath"
 "testing"
 "time"
)

// Exercises the real case through the public HTTP contract used by an
// operator screen, including visibility rules and phase restrictions.
func TestExpress1928OperatorAPIContract(t *testing.T){
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 t.Setenv("GAME_CASE_DIR",filepath.Join("..","..","cases"))
 handler,cookie,session:=testSession(t)
 check:=func(method,path string,body []byte,want int)[]byte{
  t.Helper()
  r:=requestWithCookie(handler,cookie,method,path,body)
  if r.Code!=want {t.Fatalf("%s %s returned %d, want %d: %s",method,path,r.Code,want,r.Body.String())}
  return r.Body.Bytes()
 }
 check(http.MethodGet,"/state",nil,409)
 check(http.MethodPost,"/start?caseID=case_express_1928",nil,200)
 if session.State==nil || session.State.Case.ID!="case_express_1928" {t.Fatal("case did not start")}
 defer session.State.Clock.StopTime()
 var initial gameStateView
 if err:=json.Unmarshal(check(http.MethodGet,"/state",nil,200),&initial);err!=nil {t.Fatal(err)}
 if !initial.Started || initial.Scene==nil || initial.CaseID!="case_express_1928" {t.Fatal("incomplete initial state")}
 if len(initial.FoundFactIDs)!=0 || len(initial.FoundMotiveIDs)!=0 {t.Fatal("unearned evidence exposed in initial state")}
 check(http.MethodGet,"/scene",nil,200)
 check(http.MethodPost,"/scene/choose?choiceID=victor",nil,200)
 check(http.MethodPost,"/scene/choose?choiceID=ribbon",nil,200)
 var progress gameStateView
 if err:=json.Unmarshal(check(http.MethodGet,"/state",nil,200),&progress);err!=nil {t.Fatal(err)}
 hasRibbon:=false
 for _,id:=range progress.FoundFactIDs {if id==1 {hasRibbon=true}}
 if !hasRibbon {t.Fatal("new clue missing from operator state")}
 if len(progress.FoundMotiveIDs)!=0 {t.Fatal("hidden motive exposed")}
 check(http.MethodPost,"/deduce",[]byte(`{"fact_ids":[1]}`),400)
 check(http.MethodPost,"/deduce",[]byte(`{"fact_ids":[1,999]}`),200)
 if session.State.FoundFacts[10] {t.Fatal("invalid deduction unlocked final proof")}
 check(http.MethodPost,"/accuse",[]byte(`{"suspect_id":1,"motive_id":1,"fact_ids":[10,11,12]}`),409)
 check(http.MethodPost,"/scene/choose?choiceID=back",nil,200)
 check(http.MethodPost,"/scene/choose?choiceID=salon",nil,200)
 check(http.MethodPost,"/conversation/start?id=boris_intro",nil,200)
 var talking gameStateView
 if err:=json.Unmarshal(check(http.MethodGet,"/state",nil,200),&talking);err!=nil {t.Fatal(err)}
 if talking.Conversation==nil || !talking.Paused {t.Fatal("active dialogue not reflected in state")}
 check(http.MethodPost,"/scene/choose?choiceID=back",nil,409)
 check(http.MethodPost,"/conversation/choose?choiceID=leave",nil,200)
 var after gameStateView
 if err:=json.Unmarshal(check(http.MethodGet,"/state",nil,200),&after);err!=nil {t.Fatal(err)}
 if after.Conversation!=nil || after.Paused {t.Fatal("dialogue cleanup not reflected in state")}
 session.State.InvestigationDeadline=time.Now().Add(-time.Second)
 session.State.AccusationDeadline=time.Now().Add(5*time.Minute)
 check(http.MethodPost,"/scene/choose?choiceID=back",nil,409)
 check(http.MethodPost,"/accuse",[]byte(`{"suspect_id":1,"motive_id":1,"fact_ids":[10,11,12]}`),409)
 if session.State.Ending!=nil {t.Fatal("unearned accusation finished the game")}
}

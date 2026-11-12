package server

import (
 "encoding/json"
 "net/http"
 "testing"
 "detective/internal/game"
)

func TestPrologueTravelDoesNotRevealInvestigationScene(t *testing.T) {
 handler,cookie,session:=testSession(t)
 session.Mu.Lock()
 session.State=&game.GameState{NarrativeAct:game.ActPrologue,CurrentScene:"corridor",Scenes:[]game.Scene{
  {ID:"corridor",Title:"Secret murder scene",Description:"SPOILER"},
  {ID:"salon",Title:"Secret investigation",Description:"SPOILER"},
  {ID:"compartment",Title:"Murder room",Description:"SPOILER"},
 }}
 session.Mu.Unlock()
 response:=requestWithCookie(handler,cookie,http.MethodGet,"/state",nil)
 if response.Code!=200 {t.Fatalf("state %d: %s",response.Code,response.Body.String())}
 var view gameStateView
 if err:=json.Unmarshal(response.Body.Bytes(),&view);err!=nil {t.Fatal(err)}
 if view.Scene!=nil {t.Fatal("prologue leaked investigation scene")}
 if len(view.PrologueDestinations)!=1 || view.PrologueDestinations[0].ID!="salon" {t.Fatalf("destinations: %+v",view.PrologueDestinations)}
 response=requestWithCookie(handler,cookie,http.MethodPost,"/prologue/travel?destination=compartment",nil)
 if response.Code!=409 {t.Fatalf("restricted room: %d",response.Code)}
 response=requestWithCookie(handler,cookie,http.MethodPost,"/prologue/travel?destination=salon",nil)
 if response.Code!=200 {t.Fatalf("allowed travel: %d: %s",response.Code,response.Body.String())}
 session.Mu.RLock()
 defer session.Mu.RUnlock()
 if session.State.CurrentScene!="salon" {t.Fatalf("scene=%q",session.State.CurrentScene)}
}

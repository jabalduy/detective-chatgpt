package server

import (
 "detective/internal/game"
 "net/http"
 "testing"
 "time"
)

func TestLegacyDialogueFailureDoesNotMarkAsked(t *testing.T){
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 key:=game.DialogueID{SusID:7,DialID:2}
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"test"},Clock:game.NewClock(start),
  Suspects:[]game.Suspect{{SusID:7}},
  Dialogues:[]game.Dialogue{{SusID:7,DialID:2,Action:game.Action{
   Requirements:[]game.Requirement{{Type:"clue_found",ID:9}},
   Effects:[]game.Effect{{Type:"discover_fact",ID:11}},
  }}},
  OpenDialogues:map[game.DialogueID]bool{key:true},
  AskedDialogues:map[game.DialogueID]bool{},
  FoundClues:map[int]bool{},FoundFacts:map[int]bool{},
 }
 defer session.State.Clock.StopTime()
 path:="/ask-dialogue?susID=7&dialID=2"
 r:=requestWithCookie(handler,cookie,http.MethodPost,path,nil)
 if r.Code!=http.StatusConflict{t.Fatalf("locked dialogue: %d %s",r.Code,r.Body.String())}
 if session.State.AskedDialogues[key] || session.State.FoundFacts[11]{
  t.Fatal("rejected dialogue was marked asked or granted evidence")
 }
 session.State.FoundClues[9]=true
 r=requestWithCookie(handler,cookie,http.MethodPost,path,nil)
 if r.Code!=http.StatusOK{t.Fatalf("unlocked dialogue: %d %s",r.Code,r.Body.String())}
 if !session.State.AskedDialogues[key] || !session.State.FoundFacts[11]{
  t.Fatal("successful dialogue not recorded")
 }
 session.State.Ending=&game.GameEnding{Type:"win"}
 r=requestWithCookie(handler,cookie,http.MethodPost,path,nil)
 if r.Code!=http.StatusConflict{t.Fatalf("repeated dialogue after ending: %d %s",r.Code,r.Body.String())}
 session.State.Ending=nil
 eventID:=4
 session.State.ActiveEvent=&eventID
 r=requestWithCookie(handler,cookie,http.MethodPost,path,nil)
 if r.Code!=http.StatusConflict{t.Fatalf("repeated dialogue during event: %d %s",r.Code,r.Body.String())}
}

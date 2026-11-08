package server

import (
 "detective/internal/game"
 "net/http"
 "os"
 "testing"
 "time"
)

func TestAutosaveOnSuccessfulCommand(t *testing.T){
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 start:=time.Date(2026,1,1,20,0,0,0,time.UTC)
 session.State=&game.GameState{
  Case:game.CaseInfo{ID:"test-case"},
  Clock:game.NewClock(start),CurrentScene:"corridor",
  Scenes:[]game.Scene{{ID:"corridor",Choices:[]game.Choice{{ID:"inspect",Action:game.Action{Effects:[]game.Effect{{Type:"discover_clue",ID:7}}}}}}},
  FoundClues:make(map[int]bool),OpenDialogues:make(map[game.DialogueID]bool),
 }
 defer session.State.Clock.StopTime()
 r:=requestWithCookie(handler,cookie,http.MethodPost,"/scene/choose?choiceID=inspect",nil)
 if r.Code!=200{t.Fatalf("command: %d %s",r.Code,r.Body.String())}
 saved,err:=readSave(savePath(cookie.Value))
 if err!=nil{t.Fatal(err)}
 if !saved.FoundClues[7]{t.Fatal("successful command was not saved")}
 info,err:=os.Stat(savePath(cookie.Value));if err!=nil{t.Fatal(err)}
 r=requestWithCookie(handler,cookie,http.MethodPost,"/scene/choose?choiceID=unknown",nil)
 if r.Code!=409{t.Fatalf("invalid command: %d",r.Code)}
 after,err:=os.Stat(savePath(cookie.Value));if err!=nil{t.Fatal(err)}
 if !after.ModTime().Equal(info.ModTime()){t.Fatal("failed command unexpectedly rewrote save")}
}

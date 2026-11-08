package server

import (
 "detective/internal/game"
 "net/http"
 "os"
 "path/filepath"
 "testing"
 "time"
)

func TestSaveRestoreRoundTrip(t *testing.T){
 start:=time.Date(2026,1,1,23,50,0,0,time.UTC)
 original:=&game.GameState{
  Case:game.CaseInfo{ID:"test-case"},Clock:game.NewClock(start),
  CurrentScene:"corridor",Characters:map[int]game.CharacterState{1:{Location:"corridor",Present:true,Trust:3}},
  FoundClues:map[int]bool{5:true},FoundFacts:map[int]bool{9:true},
  AskedDialogues:map[game.DialogueID]bool{{SusID:1,DialID:2}:true},
  ExaminationStarted:map[string]time.Time{"lab":start},ExaminationCompleted:map[string]bool{"lab":true},
  TriggeredEvents:map[int]bool{8:true},
 }
 original.Clock.Pause()
 path:=filepath.Join(t.TempDir(),"save.json")
 if err:=writeSave(path,snapshot(original));err!=nil{t.Fatal(err)}
 saved,err:=readSave(path);if err!=nil{t.Fatal(err)}
 restored:=&game.GameState{Clock:game.NewClock(start)}
 applySave(restored,saved)
 defer restored.Clock.StopTime()
 if restored.CurrentScene!="corridor"||restored.Characters[1].Trust!=3||!restored.FoundClues[5]||!restored.FoundFacts[9]{
  t.Fatal("basic state lost")
 }
 if !restored.AskedDialogues[game.DialogueID{SusID:1,DialID:2}]||!restored.ExaminationCompleted["lab"]||!restored.TriggeredEvents[8]{
  t.Fatal("progress state lost")
 }
 if !restored.Clock.IsPaused()||!restored.Clock.GetTime().Equal(start){t.Fatal("clock state lost")}
 original.Clock.StopTime()
}

func TestSaveEndpointAndMissingLoad(t *testing.T){
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 result:=requestWithCookie(handler,cookie,http.MethodPost,"/save",nil)
 if result.Code!=409{t.Fatalf("save before start: %d",result.Code)}
 result=requestWithCookie(handler,cookie,http.MethodPost,"/load",nil)
 if result.Code!=404{t.Fatalf("load without save: %d",result.Code)}
 session.State=&game.GameState{Case:game.CaseInfo{ID:"case_train"},Clock:game.NewClock(time.Now())}
 defer session.State.Clock.StopTime()
 result=requestWithCookie(handler,cookie,http.MethodPost,"/save",nil)
 if result.Code!=200{t.Fatalf("save: %d %s",result.Code,result.Body.String())}
 if _,err:=os.Stat(savePath(cookie.Value));err!=nil{t.Fatal(err)}
}

func TestSaveRejectsUnsupportedVersion(t *testing.T){
 path:=filepath.Join(t.TempDir(),"bad.json")
 if err:=os.WriteFile(path,[]byte(`{"version":999,"case_id":"case_train","clock_time":"2026-01-01T00:00:00Z"}`),0600);err!=nil{t.Fatal(err)}
 if _,err:=readSave(path);err==nil{t.Fatal("unsupported version accepted")}
}

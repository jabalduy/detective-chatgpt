package server

import (
 "net/http"
 "os"
 "path/filepath"
 "testing"
)

// Exercise the same path used after a real process restart: no in-memory
// session exists, but the browser retains its cookie and the save file exists.
func TestSessionRestoresAfterMemoryReset(t *testing.T){
 oldWD,err:=os.Getwd();if err!=nil{t.Fatal(err)}
 root,err:=filepath.Abs(filepath.Join(oldWD,"../.."));if err!=nil{t.Fatal(err)}
 if err:=os.Chdir(root);err!=nil{t.Fatal(err)}
 defer os.Chdir(oldWD)
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 response:=requestWithCookie(handler,cookie,http.MethodPost,"/start?caseID=case_train",nil)
 if response.Code!=200{t.Fatalf("start: %d %s",response.Code,response.Body.String())}
 session.Mu.Lock()
 session.State.FoundClues[12345]=true
 session.State.CurrentScene="recovered-scene"
 session.Mu.Unlock()
 response=requestWithCookie(handler,cookie,http.MethodPost,"/save",nil)
 if response.Code!=200{t.Fatalf("save: %d %s",response.Code,response.Body.String())}
 session.Mu.Lock()
 session.State.Clock.StopTime()
 session.Mu.Unlock()
 gamesMu.Lock()
 delete(games,cookie.Value)
 gamesMu.Unlock()
 // Simulates a new server process receiving the same browser cookie.
 response=requestWithCookie(handler,cookie,http.MethodGet,"/time",nil)
 if response.Code!=200{t.Fatalf("restored session: %d %s",response.Code,response.Body.String())}
 gamesMu.Lock()
 restored:=games[cookie.Value]
 gamesMu.Unlock()
 if restored==nil{t.Fatal("session was not reconstructed")}
 restored.Mu.RLock()
 defer restored.Mu.RUnlock()
 defer restored.State.Clock.StopTime()
 if !restored.State.FoundClues[12345] || restored.State.CurrentScene!="recovered-scene"{
  t.Fatal("progress was lost after memory reset")
 }
}

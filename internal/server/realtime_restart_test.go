package server

import (
 "net/http"
 "os"
 "path/filepath"
 "testing"
 "time"
)

// A browser's first request after a restart must see the expired ending,
// not a transient 204 response from an unrestored session.
func TestExpiredAccusationSettlesOnFirstRequestAfterRestart(t *testing.T){
 oldWD,err:=os.Getwd();if err!=nil{t.Fatal(err)}
 root,err:=filepath.Abs(filepath.Join(oldWD,"../.."));if err!=nil{t.Fatal(err)}
 if err:=os.Chdir(root);err!=nil{t.Fatal(err)}
 defer os.Chdir(oldWD)
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 response:=requestWithCookie(handler,cookie,http.MethodPost,"/start?caseID=case_train",nil)
 if response.Code!=200 {t.Fatalf("start: %d %s",response.Code,response.Body.String())}
 session.Mu.Lock()
 session.State.InvestigationDeadline=time.Now().Add(-10*time.Minute)
 session.State.AccusationDeadline=time.Now().Add(-5*time.Minute)
 session.Mu.Unlock()
 // Write directly to simulate a saved game that expired while offline.
 session.Mu.RLock()
 err=writeSave(savePath(cookie.Value),snapshot(session.State))
 session.Mu.RUnlock()
 if err!=nil {t.Fatal(err)}
 session.Mu.Lock()
 session.State.Clock.StopTime()
 session.Mu.Unlock()
 gamesMu.Lock()
 delete(games,cookie.Value)
 gamesMu.Unlock()
 response=requestWithCookie(handler,cookie,http.MethodGet,"/ending",nil)
 if response.Code!=200 {t.Fatalf("first request after restart: %d %s",response.Code,response.Body.String())}
 gamesMu.Lock()
 restored:=games[cookie.Value]
 gamesMu.Unlock()
 if restored==nil {t.Fatal("session not restored")}
 defer restored.State.Clock.StopTime()
 if restored.State.Ending==nil || restored.State.Ending.Type!="unsolved" {t.Fatal("deadline not settled after restore")}
 saved,err:=readSave(savePath(cookie.Value))
 if err!=nil || saved.Ending==nil {t.Fatalf("expired ending not saved: %v",err)}
}

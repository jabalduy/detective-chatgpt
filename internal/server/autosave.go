package server

import (
 "log"
 "os"
 "net/http"
 "strings"
)

// saveAfterMutation persists successful game commands without exposing a
// save/load button to players. The session lock is acquired only after the
// handler has released it.
func saveAfterMutation(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if !shouldAutosave(r) {
   next.ServeHTTP(w,r);return
  }
  recorder:=&statusWriter{ResponseWriter:w}
  next.ServeHTTP(recorder,r)
  if recorder.status<200 || recorder.status>=300{return}
  cookie,err:=r.Cookie("sessionID")
  if err!=nil {
   for _,raw:=range w.Header().Values("Set-Cookie") {
    if strings.HasPrefix(raw,"sessionID=") {
     cookie=&http.Cookie{Value:strings.SplitN(strings.TrimPrefix(raw,"sessionID="),";",2)[0]}
     break
    }
   }
   if cookie==nil{return}
  }
  if len(cookie.Value)!=32 || strings.ContainsAny(cookie.Value,"/\\"){return}
  gamesMu.Lock()
  session:=games[cookie.Value]
  gamesMu.Unlock()
  if session==nil{return}
  session.Mu.RLock()
  defer session.Mu.RUnlock()
  if session.State==nil || session.State.Clock==nil{return}
  if r.Method==http.MethodPost && r.URL.Path!="/start" {
   // Preserve the previous durable state as a recovery point before replacing it.
   // Do not rotate recovery on GET /state polling.
   if previous,err:=readSave(savePath(cookie.Value));err==nil {
    if err:=writeSave(recoveryPath(cookie.Value),previous);err!=nil {log.Printf("recovery checkpoint failed: %v",err)}
   } else if !os.IsNotExist(err) {log.Printf("recovery checkpoint unreadable: %v",err)}
  }
  if err:=writeSave(savePath(cookie.Value),snapshot(session.State));err!=nil{
   log.Printf("autosave session %s failed: %v",cookie.Value,err)
  }
 })
}
type statusWriter struct{
 http.ResponseWriter
 status int
}
func (w *statusWriter)WriteHeader(code int){if w.status==0{w.status=code};w.ResponseWriter.WriteHeader(code)}
func (w *statusWriter)Write(b []byte)(int,error){if w.status==0{w.status=200};return w.ResponseWriter.Write(b)}


// Some legacy GET endpoints advance scheduled events or finish examinations.
// Their resulting state changes must survive a server restart as well.
func shouldAutosave(r *http.Request) bool {
 if r.Method==http.MethodPost{return r.URL.Path!="/save" && r.URL.Path!="/load"}
 if r.Method!=http.MethodGet{return false}
 switch r.URL.Path{
 case "/event","/train","/examinations","/state":return true
 default:return false
 }
}

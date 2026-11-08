package server

import (
 "log"
 "net/http"
 "strings"
)

// saveAfterMutation persists successful game commands without exposing a
// save/load button to players. The session lock is acquired only after the
// handler has released it.
func saveAfterMutation(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Method!=http.MethodPost || r.URL.Path=="/save" || r.URL.Path=="/load" {
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

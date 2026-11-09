package server

import (
 "log"
 "net/http"
 "time"
)

// settleAccusationDeadline is request-driven: it settles an overdue session
// when the next API request arrives, even if no browser timer is running.
// It also settles immediately after a grace-period conversation finishes.
func settleAccusationDeadline(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  cookie,err:=r.Cookie("sessionID")
  if err!=nil || len(cookie.Value)!=32 {next.ServeHTTP(w,r);return}
  gamesMu.Lock()
  session:=games[cookie.Value]
  gamesMu.Unlock()
  if session==nil {next.ServeHTTP(w,r);return}
  settle:=func(){
   session.Mu.Lock()
   changed:=session.State!=nil && session.State.ExpireAccusation(time.Now())
   if changed && session.State.Clock!=nil {
    if err:=writeSave(savePath(cookie.Value),snapshot(session.State));err!=nil{
     log.Printf("save expired accusation: %v",err)
    }
   }
   session.Mu.Unlock()
  }
  settle()
  next.ServeHTTP(w,r)
  settle()
 })
}

package server

import (
 "net/http"
 "testing"
)

func TestInvestigationEndpointMethodsAndUnstartedSession(t *testing.T){
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,_:=testSession(t)
 for _,tc:=range []struct{path,allowed string}{
  {"/examinations",http.MethodGet},
  {"/examinations/start?id=lab",http.MethodPost},
  {"/evidence",http.MethodGet},
  {"/evidence/present?id=proof",http.MethodPost},
 }{
  wrong:=http.MethodPost
  if tc.allowed==http.MethodPost{wrong=http.MethodGet}
  r:=requestWithCookie(handler,cookie,wrong,tc.path,nil)
  if r.Code!=http.StatusMethodNotAllowed{
   t.Errorf("%s %s: got %d want 405",wrong,tc.path,r.Code)
  }
  r=requestWithCookie(handler,cookie,tc.allowed,tc.path,nil)
  if r.Code!=http.StatusConflict{
   t.Errorf("%s %s without game: got %d want 409",tc.allowed,tc.path,r.Code)
  }
 }
}

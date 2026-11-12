package server

import (
 "encoding/json"
 "net/http"
 "testing"
 "time"

 "detective/internal/game"
)

func TestStatePollingAutomaticallyBeginsInvestigationAndAutosaves(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 departure:=time.Date(1928,10,1,21,0,0,0,time.UTC)
 clock:=game.NewClock(departure)
 defer clock.StopTime()
 session.Mu.Lock()
 session.State=&game.GameState{
  NarrativeAct:game.ActPrologue,
  Case:game.CaseInfo{ID:"test"},
  Clock:clock,
  RouteStart:departure,
 }
 session.Mu.Unlock()

 clock.Add(77*time.Minute)
 response:=requestWithCookie(handler,cookie,http.MethodGet,"/state",nil)
 if response.Code!=http.StatusOK {t.Fatalf("status=%d body=%s",response.Code,response.Body.String())}
 var view gameStateView
 if err:=json.Unmarshal(response.Body.Bytes(),&view);err!=nil {t.Fatal(err)}
 if view.NarrativeAct!=game.ActInvestigation || view.Phase!=game.PhaseInvestigation {
  t.Fatalf("unexpected transition: act=%s phase=%s",view.NarrativeAct,view.Phase)
 }
 session.Mu.RLock()
 deadline:=session.State.InvestigationDeadline
 session.Mu.RUnlock()
 if deadline.IsZero() || time.Until(deadline)<59*time.Minute {
  t.Fatalf("investigation timer not started: %v",deadline)
 }
 saved,err:=readSave(savePath(cookie.Value))
 if err!=nil {t.Fatal(err)}
 if saved.NarrativeAct!=game.ActInvestigation || !saved.InvestigationDeadline.Equal(deadline) {
  t.Fatalf("transition not persisted: act=%s deadline=%v",saved.NarrativeAct,saved.InvestigationDeadline)
 }
 response=requestWithCookie(handler,cookie,http.MethodGet,"/state",nil)
 if response.Code!=http.StatusOK {t.Fatalf("second poll status=%d",response.Code)}
 session.Mu.RLock()
 defer session.Mu.RUnlock()
 if !session.State.InvestigationDeadline.Equal(deadline) {t.Fatal("poll reset investigation timer")}
}

func TestStatePollingDefersDiscoveryUntilActiveEventEnds(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 departure:=time.Date(1928,10,1,21,0,0,0,time.UTC)
 clock:=game.NewClock(departure)
 defer clock.StopTime()
 eventID:=42
 session.Mu.Lock()
 session.State=&game.GameState{
  NarrativeAct:game.ActPrologue,
  Case:game.CaseInfo{ID:"test"},
  Clock:clock,
  RouteStart:departure,
  ActiveEvent:&eventID,
  Events:[]game.Event{{ID:eventID,Title:"Urgent choice"}},
 }
 session.Mu.Unlock()
 clock.Add(77*time.Minute)
 response:=requestWithCookie(handler,cookie,http.MethodGet,"/state",nil)
 if response.Code!=http.StatusOK {t.Fatalf("active event poll status=%d: %s",response.Code,response.Body.String())}
 session.Mu.RLock()
 if session.State.CurrentNarrativeAct()!=game.ActPrologue || !session.State.InvestigationDeadline.IsZero() {
  session.Mu.RUnlock()
  t.Fatal("active event was interrupted by automatic discovery")
 }
 session.Mu.RUnlock()
 session.Mu.Lock()
 session.State.ActiveEvent=nil
 session.Mu.Unlock()
 response=requestWithCookie(handler,cookie,http.MethodGet,"/state",nil)
 if response.Code!=http.StatusOK {t.Fatalf("post-event poll status=%d: %s",response.Code,response.Body.String())}
 session.Mu.RLock()
 defer session.Mu.RUnlock()
 if session.State.CurrentNarrativeAct()!=game.ActInvestigation || session.State.InvestigationDeadline.IsZero() {
  t.Fatal("investigation did not start after active event ended")
 }
}

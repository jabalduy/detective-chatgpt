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

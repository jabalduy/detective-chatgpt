package server

import (
 "encoding/json"
 "net/http"
 "testing"
 "time"

 "detective/internal/game"
)

func TestDiscoveryEndpointIsGuardedAndIdempotent(t *testing.T) {
 t.Setenv("GAME_SAVE_DIR",t.TempDir())
 handler,cookie,session:=testSession(t)
 r:=requestWithCookie(handler,cookie,http.MethodPost,"/start?caseID=case_express_1928",nil)
 if r.Code!=200 {t.Fatalf("start: %d %s",r.Code,r.Body.String())}
 defer session.State.Clock.StopTime()
 r=requestWithCookie(handler,cookie,http.MethodPost,"/discovery/start",nil)
 if r.Code!=409 {t.Fatalf("legacy case should not transition: %d",r.Code)}
 session.Mu.Lock()
 session.State.NarrativeAct=game.ActPrologue
 session.State.InvestigationDeadline=time.Time{}
 session.State.AccusationDeadline=time.Time{}
 session.Mu.Unlock()
 r=requestWithCookie(handler,cookie,http.MethodPost,"/discovery/start",nil)
 if r.Code!=409 {t.Fatalf("early discovery accepted: %d",r.Code)}
 session.Mu.Lock()
 session.State.Clock.Add(77*time.Minute)
 session.Mu.Unlock()
 r=requestWithCookie(handler,cookie,http.MethodPost,"/discovery/start",nil)
 if r.Code!=200 {t.Fatalf("discovery: %d %s",r.Code,r.Body.String())}
 session.Mu.RLock()
 deadline:=session.State.InvestigationDeadline
 if session.State.CurrentNarrativeAct()!=game.ActInvestigation || !session.State.DiscoveryCompleted {session.Mu.RUnlock();t.Fatal("transition missing")}
 session.Mu.RUnlock()
 r=requestWithCookie(handler,cookie,http.MethodPost,"/discovery/start",nil)
 if r.Code!=409 {t.Fatalf("duplicate discovery: %d",r.Code)}
 session.Mu.RLock()
 if !session.State.InvestigationDeadline.Equal(deadline) {session.Mu.RUnlock();t.Fatal("deadline reset")}
 session.Mu.RUnlock()
 r=requestWithCookie(handler,cookie,http.MethodGet,"/state",nil)
 if r.Code!=200 {t.Fatalf("state: %d",r.Code)}
 var state gameStateView
 if err:=json.Unmarshal(r.Body.Bytes(),&state);err!=nil {t.Fatal(err)}
 if state.NarrativeAct!=game.ActInvestigation {t.Fatalf("state act=%s",state.NarrativeAct)}
}

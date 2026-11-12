package game

// NarrativeAct describes story progression, independently of the real-time
// investigation/accusation phase. Empty means a legacy investigation session.
type NarrativeAct string

const (
 ActPrologue NarrativeAct = "prologue"
 ActInvestigation NarrativeAct = "investigation"
 ActFinished NarrativeAct = "finished"
)

// CurrentNarrativeAct preserves compatibility with existing saves and cases.
// A prologue is opt-in until the discovery transition is implemented.
func (s *GameState) CurrentNarrativeAct() NarrativeAct {
 if s.Ending != nil { return ActFinished }
 if s.NarrativeAct == "" { return ActInvestigation }
 return s.NarrativeAct
}

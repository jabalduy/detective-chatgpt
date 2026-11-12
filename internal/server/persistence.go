package server

import (
 "detective/internal/game"
 "encoding/json"
 "errors"
 "fmt"
 "os"
 "path/filepath"
 "time"
)

const saveVersion=1

type saveData struct {
 Version int `json:"version"`
 NarrativeAct game.NarrativeAct `json:"narrative_act,omitempty"`
 DiscoveryCompleted bool `json:"discovery_completed,omitempty"`
 UsedHints map[string]bool `json:"used_hints,omitempty"`
 HintCount int `json:"hint_count,omitempty"`
 AskedQuestionsBySuspect map[int]int `json:"asked_questions_by_suspect,omitempty"`
 CaseID string `json:"case_id"`
 InvestigationDeadline time.Time `json:"investigation_deadline,omitempty"`
 AccusationDeadline time.Time `json:"accusation_deadline,omitempty"`
 ClockTime time.Time `json:"clock_time"`
 ClockPaused bool `json:"clock_paused"`
 CurrentScene string `json:"current_scene"`
 Characters map[int]game.CharacterState `json:"characters"`
 ActiveConversation string `json:"active_conversation"`
 ConversationNode string `json:"conversation_node"`
 ActiveEvent *int `json:"active_event"`
 EventDeadline time.Time `json:"event_deadline"`
 EventRemainingSeconds float64 `json:"event_remaining_seconds,omitempty"`
 ActiveDialogue *game.DialogueID `json:"active_dialogue"`
 FoundClues map[int]bool `json:"found_clues"`
 FoundFacts map[int]bool `json:"found_facts"`
 FoundMotives map[int]bool `json:"found_motives"`
 SearchedObjects map[int]bool `json:"searched_objects"`
 SolvedDeductions map[int]bool `json:"solved_deductions"`
 OpenLocations map[int]bool `json:"open_locations"`
 OpenObjects map[int]bool `json:"open_objects"`
 AskedDialogues []game.DialogueID `json:"asked_dialogues"`
 OpenDialogues []game.DialogueID `json:"open_dialogues"`
 ConfrontedDialogues []game.DialogueID `json:"confronted_dialogues"`
 ProcessedStations map[string]bool `json:"processed_stations"`
 TriggeredEvents map[int]bool `json:"triggered_events"`
 ExaminationStarted map[string]time.Time `json:"examination_started"`
 ExaminationCompleted map[string]bool `json:"examination_completed"`
 PresentedEvidence map[string]bool `json:"presented_evidence"`
 Ending *game.GameEnding `json:"ending"`
}

func dialogueKeys(m map[game.DialogueID]bool)[]game.DialogueID{
 out:=make([]game.DialogueID,0,len(m))
 for k,v:=range m {if v{out=append(out,k)}}
 return out
}
func dialogueMap(keys []game.DialogueID)map[game.DialogueID]bool{
 out:=make(map[game.DialogueID]bool,len(keys))
 for _,k:=range keys{out[k]=true}
 return out
}
func snapshot(s *game.GameState)saveData{
 remaining:=0.0
 if s.ActiveEvent!=nil && !s.EventDeadline.IsZero(){remaining=time.Until(s.EventDeadline).Seconds();if remaining<0{remaining=0}}
 return saveData{
  Version:saveVersion,NarrativeAct:s.NarrativeAct,DiscoveryCompleted:s.DiscoveryCompleted,UsedHints:s.UsedHints,HintCount:s.HintCount,AskedQuestionsBySuspect:s.AskedQuestionsBySuspect,CaseID:s.Case.ID,InvestigationDeadline:s.InvestigationDeadline,AccusationDeadline:s.AccusationDeadline,ClockTime:s.Clock.GetTime(),ClockPaused:s.Clock.IsPaused(),
  CurrentScene:s.CurrentScene,Characters:s.Characters,ActiveConversation:s.ActiveConversation,ConversationNode:s.ConversationNode,
  ActiveEvent:s.ActiveEvent,EventDeadline:s.EventDeadline,EventRemainingSeconds:remaining,ActiveDialogue:s.ActiveDialogue,
  FoundClues:s.FoundClues,FoundFacts:s.FoundFacts,FoundMotives:s.FoundMotives,
  SearchedObjects:s.SearchedObjects,SolvedDeductions:s.SolvedDeductions,OpenLocations:s.OpenLocations,OpenObjects:s.OpenObjects,
  AskedDialogues:dialogueKeys(s.AskedDialogues),OpenDialogues:dialogueKeys(s.OpenDialogues),ConfrontedDialogues:dialogueKeys(s.ConfrontedDialogues),
  ProcessedStations:s.ProcessedStations,TriggeredEvents:s.TriggeredEvents,ExaminationStarted:s.ExaminationStarted,
  ExaminationCompleted:s.ExaminationCompleted,PresentedEvidence:s.PresentedEvidence,Ending:s.Ending,
 }
}
func applySave(s *game.GameState,d saveData){
 s.Clock.StopTime()
 s.Clock=game.NewClock(d.ClockTime)
 s.NarrativeAct=d.NarrativeAct
 s.DiscoveryCompleted=d.DiscoveryCompleted
 s.UsedHints=d.UsedHints
 s.HintCount=d.HintCount
 s.AskedQuestionsBySuspect=d.AskedQuestionsBySuspect
 if s.AskedQuestionsBySuspect==nil{s.AskedQuestionsBySuspect=make(map[int]int)}
 if s.UsedHints==nil{s.UsedHints=make(map[string]bool)}
 s.InvestigationDeadline=d.InvestigationDeadline
 s.AccusationDeadline=d.AccusationDeadline
 s.CurrentScene=d.CurrentScene
 s.Characters=d.Characters
 s.ActiveConversation=d.ActiveConversation;s.ConversationNode=d.ConversationNode
 s.ActiveEvent=d.ActiveEvent;s.EventDeadline=d.EventDeadline;s.ActiveDialogue=d.ActiveDialogue
 // Urgent decisions use wall-clock deadlines. A server restart must not
 // replenish the countdown; overdue events resolve on the next event tick.
 // EventRemainingSeconds is retained only for compatibility with old saves.
 s.FoundClues=d.FoundClues;s.FoundFacts=d.FoundFacts;s.FoundMotives=d.FoundMotives
 s.SearchedObjects=d.SearchedObjects;s.SolvedDeductions=d.SolvedDeductions
 s.OpenLocations=d.OpenLocations;s.OpenObjects=d.OpenObjects
 s.AskedDialogues=dialogueMap(d.AskedDialogues)
 s.OpenDialogues=dialogueMap(d.OpenDialogues)
 s.ConfrontedDialogues=dialogueMap(d.ConfrontedDialogues)
 s.ProcessedStations=d.ProcessedStations;s.TriggeredEvents=d.TriggeredEvents
 s.ExaminationStarted=d.ExaminationStarted;s.ExaminationCompleted=d.ExaminationCompleted
 s.PresentedEvidence=d.PresentedEvidence;s.Ending=d.Ending
 // Older or minimal saves may omit empty maps. Restore writable maps so
 // subsequent evidence and exploration effects cannot panic.
 if s.FoundClues==nil{s.FoundClues=make(map[int]bool)}
 if s.FoundFacts==nil{s.FoundFacts=make(map[int]bool)}
 if s.FoundMotives==nil{s.FoundMotives=make(map[int]bool)}
 if s.SearchedObjects==nil{s.SearchedObjects=make(map[int]bool)}
 if s.SolvedDeductions==nil{s.SolvedDeductions=make(map[int]bool)}
 if s.OpenLocations==nil{s.OpenLocations=make(map[int]bool)}
 if s.OpenObjects==nil{s.OpenObjects=make(map[int]bool)}
 if s.ProcessedStations==nil{s.ProcessedStations=make(map[string]bool)}
 if s.TriggeredEvents==nil{s.TriggeredEvents=make(map[int]bool)}
 if s.ExaminationStarted==nil{s.ExaminationStarted=make(map[string]time.Time)}
 if s.ExaminationCompleted==nil{s.ExaminationCompleted=make(map[string]bool)}
 if s.PresentedEvidence==nil{s.PresentedEvidence=make(map[string]bool)}
 if s.Characters==nil{s.Characters=make(map[int]game.CharacterState)}
 if d.ClockPaused || d.ActiveConversation!="" || d.ActiveDialogue!=nil || d.Ending!=nil {s.Clock.Pause()}
 go s.Clock.RunTime()
}
func savePath(sessionID string)string{return filepath.Join(saveDirectory(),sessionID+".json")}
func saveDirectory()string{
 if v:=os.Getenv("GAME_SAVE_DIR");v!=""{return v}
 return filepath.Join("data","saves")
}
func writeSave(path string,d saveData)error{
 data,err:=json.MarshalIndent(d,"","  ");if err!=nil{return err}
 if err=os.MkdirAll(filepath.Dir(path),0700);err!=nil{return err}
 f,err:=os.CreateTemp(filepath.Dir(path),".save-*.tmp");if err!=nil{return err}
 defer os.Remove(f.Name())
 if err=f.Chmod(0600);err!=nil{f.Close();return err}
 if _,err=f.Write(data);err!=nil{f.Close();return err}
 if err=f.Sync();err!=nil{f.Close();return err}
 if err=f.Close();err!=nil{return err}
 return os.Rename(f.Name(),path)
}
func readSave(path string)(saveData,error){
 var d saveData
 data,err:=os.ReadFile(path);if err!=nil{return d,err}
 if err=json.Unmarshal(data,&d);err!=nil{return d,err}
 if d.Version!=saveVersion{return d,fmt.Errorf("unsupported save version %d",d.Version)}
 if d.CaseID=="" || d.ClockTime.IsZero(){return d,errors.New("incomplete save")}
 return d,nil
}

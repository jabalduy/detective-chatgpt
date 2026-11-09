package game

import "time"

type GameState struct {
 Hints []Hint
 UsedHints map[string]bool
 HintCount int
 AskedQuestionsBySuspect map[int]int
 // Real-world deadlines are independent of the in-world train clock.
 InvestigationDeadline time.Time
 AccusationDeadline time.Time
 EndingRules []EndingRule
 Ending *GameEnding
 Examinations []Examination
 ExaminationStarted map[string]time.Time
 ExaminationCompleted map[string]bool
 EvidencePresentations []EvidencePresentation
 PresentedEvidence map[string]bool
	Stations []Station
	RouteStart time.Time
	ProcessedStations map[string]bool
	TriggeredEvents map[int]bool
	EventDeadline time.Time
	Conversations []Conversation
	Characters map[int]CharacterState
	ActiveConversation string
	ConversationNode string
	Scenes []Scene
	CurrentScene string
	Case CaseInfo

	Suspects    []Suspect
	Clues       []Clue
	Locations   []Location
	Objects     []Object
	Dialogues   []Dialogue
	Facts       []Fact
	Deductions  []Deduction
	Motives     []Motive
	Events      []Event
	Clock       *Clock
	ActiveEvent *int

	ActiveDialogue *DialogueID

	FoundClues          map[int]bool
	FoundFacts          map[int]bool
	SearchedObjects     map[int]bool
	AskedDialogues      map[DialogueID]bool
	OpenDialogues       map[DialogueID]bool
	SolvedDeductions    map[int]bool
	FoundMotives        map[int]bool
	OpenLocations       map[int]bool
	OpenObjects         map[int]bool
	ConfrontedDialogues map[DialogueID]bool
}

func InitOpenDialogues(dialogues []Dialogue) map[DialogueID]bool {
	openDialogues := make(map[DialogueID]bool)

	for _, dialogue := range dialogues {
		if dialogue.InitiallyOpen {
			dialKey := DialogueID{
				SusID:  dialogue.SusID,
				DialID: dialogue.DialID,
			}

			openDialogues[dialKey] = true
		}
	}

	return openDialogues
}

func InitOpenLocations(locations []Location) map[int]bool {
	openLocations := make(map[int]bool)

	for _, location := range locations {
		if location.InitiallyOpen {
			openLocations[location.LocID] = true
		}
	}

	return openLocations
}

func InitOpenObjects(objects []Object) map[int]bool {
	openObjects := make(map[int]bool)

	for _, object := range objects {
		if object.InitiallyOpen {
			openObjects[object.ObjID] = true
		}
	}

	return openObjects
}

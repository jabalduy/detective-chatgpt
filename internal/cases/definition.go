package cases

import "detective/internal/game"

type CaseDefinition struct {
	game.CaseInfo

	Clues      []game.Clue      `json:"-"`
	Locations  []game.Location  `json:"-"`
	Objects    []game.Object    `json:"-"`
	Suspects   []game.Suspect   `json:"-"`
	Dialogues  []game.Dialogue  `json:"-"`
	Facts      []game.Fact      `json:"-"`
	Deductions []game.Deduction `json:"-"`
	Motives    []game.Motive    `json:"-"`
	Events     []game.Event     `json:"-"`
	StartTime  game.StartTime   `json:"-"`
}

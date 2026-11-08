package game

func CanConfront(state *GameState, dialogue Dialogue) bool {
	switch {
	case dialogue.Confrontation == nil:
		return false

	case !state.AskedDialogues[dialogue.ID()]:
		return false

	case state.ConfrontedDialogues[dialogue.ID()]:
		return false

	case !RequirementsMet(state, dialogue.Confrontation.Requirements):
		return false

	default:
		return true
	}
}

func PerformConfrontation(state *GameState, dialogue Dialogue) bool {
	if !CanConfront(state, dialogue) {
		return false
	}

	valid := PerformAction(state, dialogue.Confrontation.Action)
	if !valid {
		return false
	}

	state.ConfrontedDialogues[dialogue.ID()] = true

	return true
}

package game

func requirementMet(state *GameState, req Requirement) bool {
	switch req.Type {
	case "clue_found":
		return state.FoundClues[req.ID]

	case "object_searched":
		return state.SearchedObjects[req.ID]

	case "fact_found":
		return state.FoundFacts[req.ID]

	case "dialogue_asked":
		dialogueID := DialogueID{
			SusID:  req.SusID,
			DialID: req.DialID,
		}

		return state.AskedDialogues[dialogueID]

	default:
		return false
	}
}

func RequirementsMet(state *GameState, requirements []Requirement) bool {
	for _, req := range requirements {
		if !requirementMet(state, req) {
			return false
		}
	}

	return true
}

func UpdateOpenDialogues(state *GameState) int {
	unlocked := 0

	for _, dialogue := range state.Dialogues {
		id := dialogue.ID()

		if state.OpenDialogues[id] {
			continue
		}

		if len(dialogue.Requirements) == 0 {
			continue
		}

		if RequirementsMet(state, dialogue.Requirements) {
			state.OpenDialogues[id] = true
			unlocked++
		}
	}

	return unlocked
}

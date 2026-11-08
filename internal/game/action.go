package game

import "time"

func CanPerformAction(state *GameState, action Action) bool {
	return RequirementsMet(state, action.Requirements)
}

func PerformAction(state *GameState, action Action) bool {
	if !CanPerformAction(state, action) {
		return false
	}

	ApplyEffects(state, action.Effects)
	UpdateOpenDialogues(state)

	timeCost := time.Duration(action.TimeCost)
	state.Clock.Add(timeCost * time.Second)

	return true
}

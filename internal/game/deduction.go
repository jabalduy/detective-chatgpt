package game

type Deduction struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Required []int  `json:"required_facts"`

	Action
}

func CanSolveDeduction(state *GameState, deduction Deduction) bool {
	if state.SolvedDeductions[deduction.ID] {
		return false
	}

	for _, factID := range deduction.Required {
		if !state.FoundFacts[factID] {
			return false
		}
	}

	return true
}

func SolveDeduction(state *GameState, deduction Deduction) bool {
	if !CanSolveDeduction(state, deduction) {
		return false
	}

	state.SolvedDeductions[deduction.ID] = true

	ApplyEffects(state, deduction.Effects)
	UpdateOpenDialogues(state)

	return true
}

func FindDeduction(state *GameState, factIDs []int) *Deduction {
	for i := range state.Deductions {
		deduction := &state.Deductions[i]

		if sameFacts(deduction.Required, factIDs) {
			return deduction
		}
	}

	return nil
}

func sameFacts(required []int, selected []int) bool {
	if len(required) != len(selected) {
		return false
	}

	selectedSet := make(map[int]bool, len(selected))

	for _, id := range selected {
		selectedSet[id] = true
	}

	for _, id := range required {
		if !selectedSet[id] {
			return false
		}
	}

	return true
}

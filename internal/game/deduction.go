package game

type Deduction struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Required []int  `json:"required_facts"`

	Action
}

func CanSolveDeduction(state *GameState, deduction Deduction) bool {
	if state.Ending != nil || state.SolvedDeductions[deduction.ID] || !CanPerformAction(state, deduction.Action) {
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

	if !PerformAction(state, deduction.Action) {
		return false
	}
	if state.SolvedDeductions == nil {
		state.SolvedDeductions = make(map[int]bool)
	}
	state.SolvedDeductions[deduction.ID] = true
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

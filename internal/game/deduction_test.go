package game

import "testing"

func TestCanSolveDeductionCases(t *testing.T) {
	tests := []struct {
		name     string
		found    map[int]bool
		solved   map[int]bool
		required []int
		want     bool
	}{
		{
			name: "all facts found",
			found: map[int]bool{
				1: true,
				2: true,
			},
			solved:   map[int]bool{},
			required: []int{1, 2},
			want:     true,
		},
		{
			name: "missing fact",
			found: map[int]bool{
				1: true,
			},
			solved:   map[int]bool{},
			required: []int{1, 2},
			want:     false,
		},
		{
			name: "already solved",
			found: map[int]bool{
				1: true,
				2: true,
			},
			solved: map[int]bool{
				1: true,
			},
			required: []int{1, 2},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &GameState{
				FoundFacts:       tt.found,
				SolvedDeductions: tt.solved,
			}

			deduction := Deduction{
				ID:       1,
				Required: tt.required,
			}

			got := CanSolveDeduction(state, deduction)

			if got != tt.want {
				t.Errorf(
					"CanSolveDeduction() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestSolveDeductionDiscoversFact(t *testing.T) {
	state := &GameState{
		FoundFacts: map[int]bool{
			1: true,
			2: true,
		},
		SolvedDeductions: make(map[int]bool),
		OpenDialogues:    make(map[DialogueID]bool),
	}

	deduction := Deduction{
		ID:       1,
		Required: []int{1, 2},
		Action: Action{
			Effects: []Effect{
				{
					Type: "discover_fact",
					ID:   3,
				},
			},
		},
	}

	ok := SolveDeduction(state, deduction)

	if !ok {
		t.Fatal("expected deduction to be solved")
	}

	if !state.SolvedDeductions[1] {
		t.Error("expected deduction 1 to be marked as solved")
	}

	if !state.FoundFacts[3] {
		t.Error("expected fact 3 to be discovered")
	}
}

func TestSolveDeductionDiscoversMotive(t *testing.T) {
	state := &GameState{
		FoundFacts: map[int]bool{
			1: true,
			2: true,
		},
		FoundMotives:     make(map[int]bool),
		SolvedDeductions: make(map[int]bool),
		OpenDialogues:    make(map[DialogueID]bool),
	}

	deduction := Deduction{
		ID:       1,
		Required: []int{1, 2},
		Action: Action{
			Effects: []Effect{
				{
					Type: "discover_motive",
					ID:   10,
				},
			},
		},
	}

	success := SolveDeduction(state, deduction)

	if !success {
		t.Fatal("expected deduction to be solved")
	}

	if !state.FoundMotives[10] {
		t.Fatal("expected motive 10 to be discovered")
	}

	if !state.SolvedDeductions[1] {
		t.Fatal("expected deduction 1 to be marked as solved")
	}
}

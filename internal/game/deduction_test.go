package game

import ("testing"; "time")

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

func TestDeductionRequiresEvidenceAndActionConditions(t *testing.T) {
 start:=time.Date(2026,1,1,12,0,0,0,time.UTC)
 s:=&GameState{
  Clock:NewClock(start),
  FoundFacts:map[int]bool{1:true,2:true},
  FoundClues:map[int]bool{},
 }
 defer s.Clock.StopTime()
 d:=Deduction{ID:7,Required:[]int{1,2},Action:Action{
  Requirements:[]Requirement{{Type:"clue_found",ID:9}},
  TimeCost:30,
  Effects:[]Effect{{Type:"discover_fact",ID:3}},
 }}
 if CanSolveDeduction(s,d) || SolveDeduction(s,d) {t.Fatal("locked deduction accepted")}
 if s.SolvedDeductions[7] || s.FoundFacts[3] {t.Fatal("rejected deduction mutated state")}
 s.FoundClues[9]=true
 if !SolveDeduction(s,d) {t.Fatal("unlocked deduction rejected")}
 if !s.SolvedDeductions[7] || !s.FoundFacts[3] {t.Fatal("deduction consequences missing")}
 if got:=s.Clock.GetTime();got.Before(start.Add(30*time.Second)) {t.Fatalf("time cost missing: %v",got)}
 if SolveDeduction(s,d) {t.Fatal("deduction repeated")}
}

func TestDeductionRejectsDuplicateFactIDs(t *testing.T) {
 if sameFacts([]int{1,2},[]int{1,1}) {t.Fatal("duplicate facts accepted as distinct evidence")}
 if !sameFacts([]int{1,2},[]int{2,1}) {t.Fatal("order-independent evidence rejected")}
}

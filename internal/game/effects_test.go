package game

import "testing"

func TestApplyEffectsOpenLocation(t *testing.T) {
	state := &GameState{
		OpenLocations: make(map[int]bool),
	}

	effects := []Effect{
		{
			Type: "open_location",
			ID:   4,
		},
	}

	ApplyEffects(state, effects)

	if !state.OpenLocations[4] {
		t.Fatalf("location 4 should be open after open_location effect")
	}
}

package game

func applyEffect(state *GameState, effect Effect) {
	switch effect.Type {
	case "character_trust_up":
		c:=state.Characters[effect.ID];c.Trust++;state.Characters[effect.ID]=c
	case "character_trust_down":
		c:=state.Characters[effect.ID];c.Trust--;state.Characters[effect.ID]=c
	case "character_leave":
		c:=state.Characters[effect.ID];c.Present=false;state.Characters[effect.ID]=c
	case "character_appear":
		c:=state.Characters[effect.ID];c.Present=true;state.Characters[effect.ID]=c
	case "discover_clue":
		state.FoundClues[effect.ID] = true

	case "discover_fact":
		state.FoundFacts[effect.ID] = true

	case "discover_motive":
		state.FoundMotives[effect.ID] = true

	case "open_location":
		state.OpenLocations[effect.ID] = true

	case "open_object":
		state.OpenObjects[effect.ID] = true

	case "trigger_event":
		TriggerEvent(state, effect.ID)

	case "open_dialogue":
		id := DialogueID{
			SusID:  effect.SusID,
			DialID: effect.DialID,
		}

		state.OpenDialogues[id] = true
	}
}

func ApplyEffects(state *GameState, effects []Effect) {
	for _, effect := range effects {
		applyEffect(state, effect)
	}
}

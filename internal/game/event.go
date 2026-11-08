package game

type Event struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	TimeLimit   int           `json:"time_limit"`
	Options     []EventOption `json:"options"`
}

type EventOption struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Action
}

func TriggerEvent(state *GameState, eventID int) bool {
	for _, event := range state.Events {
		if event.ID == eventID {
			state.ActiveEvent = &eventID
			return true
		}
	}

	return false
}

func ChooseEventOption(state *GameState, optionID int) bool {
	if state.ActiveEvent == nil {
		return false
	}

	eventID := *state.ActiveEvent

	for _, event := range state.Events {
		if event.ID != eventID {
			continue
		}

		for _, option := range event.Options {
			if option.ID != optionID {
				continue
			}

			if !PerformAction(state, option.Action) {
				return false
			}

			state.ActiveEvent = nil
			return true
		}
		return false
	}
	return false
}

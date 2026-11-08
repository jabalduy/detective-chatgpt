package game

func KeyDial(game *GameState) int {
	keys := 0

	for _, dial := range game.Dialogues {
		dialKey := DialogueID{
			SusID:  dial.SusID,
			DialID: dial.DialID,
		}

		if dial.Key && game.AskedDialogues[dialKey] {
			keys += 1
		}
	}
	return keys
}

func KeyClues(game *GameState) int {
	keys := 0
	for _, clue := range game.Clues {
		if clue.Key && game.FoundClues[clue.ObjID] {
			keys += 1
		}
	}
	return keys
}

func KeyObj(game *GameState) int {
	keys := 0
	for _, obj := range game.Objects {
		if obj.Key && game.SearchedObjects[obj.ObjID] {
			keys += 1
		}
	}
	return keys
}

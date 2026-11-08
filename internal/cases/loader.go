package cases

import (
	"detective/internal/game"
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadCase(path string) (*CaseDefinition, error) {
	data, err := os.ReadFile(filepath.Join(path, "case.json"))
	if err != nil {
		return nil, err
	}

	var caseDef CaseDefinition

	err = json.Unmarshal(data, &caseDef)
	if err != nil {
		return nil, err
	}

	clues, err := LoadClues(filepath.Join(path, "clues.json"))
	if err != nil {
		return nil, err
	}

	locations, err := LoadLocations(filepath.Join(path, "locations.json"))
	if err != nil {
		return nil, err
	}

	objects, err := LoadObjects(filepath.Join(path, "objects.json"))
	if err != nil {
		return nil, err
	}

	suspects, err := LoadSuspects(filepath.Join(path, "suspects.json"))
	if err != nil {
		return nil, err
	}

	dialogues, err := LoadDialogues(filepath.Join(path, "dialogues.json"))
	if err != nil {
		return nil, err
	}

	facts, err := LoadFacts(filepath.Join(path, "facts.json"))
	if err != nil {
		return nil, err
	}

	deductions, err := LoadDeductions(filepath.Join(path, "deductions.json"))
	if err != nil {
		return nil, err
	}

	motives, err := LoadMotives(filepath.Join(path, "motives.json"))
	if err != nil {
		return nil, err
	}

	startTime, err := LoadStartTime(filepath.Join(path, "startTime.json"))
	if err != nil {
		return nil, err
	}

	events, err := LoadEvents(filepath.Join(path, "events.json"))
	if err != nil {
		return nil, err
	}

	caseDef.Clues = clues
	caseDef.Locations = locations
	caseDef.Objects = objects
	caseDef.Suspects = suspects
	caseDef.Dialogues = dialogues
	caseDef.Facts = facts
	caseDef.Deductions = deductions
	caseDef.Motives = motives
	caseDef.StartTime = startTime
	caseDef.Events = events

	return &caseDef, nil
}

func LoadClues(path string) ([]game.Clue, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var clues []game.Clue

	err = json.Unmarshal(data, &clues)
	if err != nil {
		return nil, err
	}

	return clues, nil
}

func LoadLocations(path string) ([]game.Location, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var locations []game.Location

	err = json.Unmarshal(data, &locations)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func LoadObjects(path string) ([]game.Object, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var objects []game.Object

	err = json.Unmarshal(data, &objects)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func LoadSuspects(path string) ([]game.Suspect, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var suspects []game.Suspect

	err = json.Unmarshal(data, &suspects)
	if err != nil {
		return nil, err
	}

	return suspects, nil
}

func LoadDialogues(path string) ([]game.Dialogue, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var dialogues []game.Dialogue

	err = json.Unmarshal(data, &dialogues)
	if err != nil {
		return nil, err
	}

	return dialogues, nil
}

func LoadFacts(path string) ([]game.Fact, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var facts []game.Fact

	err = json.Unmarshal(data, &facts)
	if err != nil {
		return nil, err
	}

	return facts, nil
}

func LoadDeductions(path string) ([]game.Deduction, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var deductions []game.Deduction

	err = json.Unmarshal(data, &deductions)
	if err != nil {
		return nil, err
	}

	return deductions, nil
}

func LoadMotives(path string) ([]game.Motive, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var motives []game.Motive

	err = json.Unmarshal(data, &motives)
	if err != nil {
		return nil, err
	}

	return motives, nil
}

func LoadStartTime(path string) (game.StartTime, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return game.StartTime{}, err
	}

	var startTime game.StartTime

	err = json.Unmarshal(data, &startTime)
	if err != nil {
		return game.StartTime{}, err
	}

	return startTime, nil
}

func LoadEvents(path string) ([]game.Event, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var events []game.Event

	err = json.Unmarshal(data, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

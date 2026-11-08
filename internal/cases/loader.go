package cases

import (
	"detective/internal/game"
	"encoding/json"
	"os"
	"path/filepath"
	"errors"
	"fmt"
	"strings"
	"time"
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
	if err != nil && !errors.Is(err, os.ErrNotExist) {
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

	scenes, err := LoadScenes(filepath.Join(path, "scenes.json"))
	if err != nil && !errors.Is(err, os.ErrNotExist) { return nil, err }
	caseDef.Scenes = scenes
	conversations, err := LoadConversations(filepath.Join(path, "conversations.json"))
	if err != nil && !errors.Is(err, os.ErrNotExist) { return nil, err }
	caseDef.Conversations = conversations
	characters, err := LoadCharacters(filepath.Join(path, "characters.json"))
	if err != nil && !errors.Is(err, os.ErrNotExist) { return nil, err }
	caseDef.Characters = characters
	stations,err:=LoadStations(filepath.Join(path,"stations.json"))
	if err!=nil && !errors.Is(err,os.ErrNotExist){return nil,err}
	caseDef.Stations=stations
 examinations,err:=LoadExaminations(filepath.Join(path,"examinations.json"))
 if err!=nil && !errors.Is(err,os.ErrNotExist){return nil,err}
 caseDef.Examinations=examinations
 presentations,err:=LoadEvidencePresentations(filepath.Join(path,"evidence_presentations.json"))
 if err!=nil && !errors.Is(err,os.ErrNotExist){return nil,err}
 caseDef.EvidencePresentations=presentations
 endingRules,err:=LoadEndingRules(filepath.Join(path,"ending_rules.json"))
 if err!=nil && !errors.Is(err,os.ErrNotExist){return nil,err}
 caseDef.EndingRules=endingRules

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
 if err != nil { return game.StartTime{}, err }
 var raw struct { StartTime string `json:"start_time"` }
 if err := json.Unmarshal(data, &raw); err != nil { return game.StartTime{}, err }
 value := strings.TrimSpace(raw.StartTime)
 for _, layout := range []string{"15:04", time.RFC3339} {
  parsed, err := time.Parse(layout, value)
  if err == nil { return game.StartTime{StartTime: parsed}, nil }
 }
 return game.StartTime{}, fmt.Errorf("invalid start_time %q: expected HH:MM or RFC3339", value)
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

func LoadScenes(path string) ([]game.Scene, error) {
 data, err := os.ReadFile(path)
 if err != nil { return nil, err }
 var scenes []game.Scene
 if err := json.Unmarshal(data, &scenes); err != nil { return nil, err }
 return scenes,nil
}

func LoadConversations(path string)([]game.Conversation,error){
 data,err:=os.ReadFile(path);if err!=nil{return nil,err}
 var items []game.Conversation
 err=json.Unmarshal(data,&items);return items,err
}
func LoadCharacters(path string)(map[int]game.CharacterState,error){
 data,err:=os.ReadFile(path);if err!=nil{return nil,err}
 var items map[int]game.CharacterState
 err=json.Unmarshal(data,&items);return items,err
}

func LoadStations(path string)([]game.Station,error){
 data,err:=os.ReadFile(path);if err!=nil{return nil,err}
 var items []game.Station
 err=json.Unmarshal(data,&items);return items,err
}

func LoadExaminations(path string)([]game.Examination,error){
 data,err:=os.ReadFile(path);if err!=nil{return nil,err}
 var items []game.Examination
 err=json.Unmarshal(data,&items);return items,err
}
func LoadEvidencePresentations(path string)([]game.EvidencePresentation,error){
 data,err:=os.ReadFile(path);if err!=nil{return nil,err}
 var items []game.EvidencePresentation
 err=json.Unmarshal(data,&items);return items,err
}

func LoadEndingRules(path string)([]game.EndingRule,error){
 data,err:=os.ReadFile(path);if err!=nil{return nil,err}
 var items []game.EndingRule
 err=json.Unmarshal(data,&items);return items,err
}

package server

import (
	"detective/internal/game"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// ЛОКАЦИИ, ОБЪЕКТЫ, РАССЛЕДОВАНИЕ
func locationsHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getGame(w, r)
	if err != nil {
		log.Printf("new game: %v", err)
		http.Error(w, "Не удалось загрузить игру", http.StatusInternalServerError)
		return
	}

	session.Mu.RLock()
	defer session.Mu.RUnlock()
	state := session.State

	if state == nil {
		http.Error(w, "Игра не запущена", http.StatusConflict)
		return
	}

	openLocations := []game.Location{}

	for _, location := range state.Locations {
		if !state.OpenLocations[location.LocID] {
			continue
		}

		openLocations = append(openLocations, location)
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(openLocations)
	if err != nil {
		log.Printf("locations: encode response: %v", err)
	}
}

func objectsHandler(w http.ResponseWriter, r *http.Request) {
	// ПОЛУЧИТЬ ИГРУ
	session, err := getGame(w, r)
	if err != nil {
		log.Printf("new state: %v", err)
		http.Error(w, "Не удалось загрузить игру", http.StatusInternalServerError)
		return
	}

	session.Mu.RLock()
	defer session.Mu.RUnlock()
	state := session.State

	if state == nil {
		http.Error(w, "Игра не запущена", http.StatusConflict)
		return
	}

	// ОБРАБРОТКА ВЫБОРА КОМНАТЫ
	locIDStr := r.URL.Query().Get("locID")

	locID, err := strconv.Atoi(locIDStr)
	if err != nil {
		http.Error(w, "Некорректный locID", http.StatusBadRequest)
		return
	}

	existLocation := false
	for _, location := range state.Locations {
		if location.LocID == locID {
			existLocation = true
			break
		}
	}

	if !existLocation {
		http.Error(
			w,
			"Нет такой локации",
			http.StatusNotFound,
		)
		return
	}

	if !state.OpenLocations[locID] {
		http.Error(
			w,
			"Недоступная локация",
			http.StatusForbidden,
		)
		return
	}

	objectsResponse := make([]game.ObjectResponse, 0)

	for _, obj := range state.Objects {
		if obj.LocID == locID &&
			state.OpenObjects[obj.ObjID] {
			newObj := game.ObjectResponse{
				LocID:    obj.LocID,
				Name:     obj.Name,
				About:    obj.About,
				ObjID:    obj.ObjID,
				Key:      obj.Key,
				Searched: state.SearchedObjects[obj.ObjID],
			}
			objectsResponse = append(objectsResponse, newObj)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(objectsResponse)
	if err != nil {
		log.Printf("objects: encode response: %v", err)
	}
}

func inspectObjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Некорректная команда", http.StatusMethodNotAllowed)
		return
	}

	session, err := getGame(w, r)
	if err != nil {
		log.Printf("new game: %v", err)
		http.Error(w, "Не удалось загрузить игру", http.StatusInternalServerError)
		return
	}

	session.Mu.Lock()
	defer session.Mu.Unlock()
	state := session.State

	if state == nil {
		http.Error(w, "Игра не запущена", http.StatusConflict)
		return
	}

	objIDstr := r.URL.Query().Get("objID")

	objID, err := strconv.Atoi(objIDstr)
	if err != nil {
		http.Error(w, "Некорректный objID", http.StatusBadRequest)
		return
	}

	for i := range state.Objects {
		if state.Objects[i].ObjID != objID {
			continue
		}

		if !state.OpenObjects[objID] {
			http.Error(
				w,
				"Объект недоступен",
				http.StatusForbidden,
			)
			return
		}

		object := state.Objects[i]

		response := InspectResponse{
			Object: state.Objects[i],
		}

		if !state.SearchedObjects[objID] {
			if !game.PerformAction(state, object.InspectAction) {
				http.Error(
					w,
					"Условия для осмотра объекта не выполнены",
					http.StatusForbidden,
				)
				return
			}

			state.SearchedObjects[objID] = true
		}

		if state.FoundClues[objID] {
			for clueIndex := range state.Clues {
				if state.Clues[clueIndex].ObjID != objID {
					continue
				}

				response.ClueFound = true
				response.ClueName = state.Clues[clueIndex].Name
				response.ClueAbout = state.Clues[clueIndex].About

				break
			}
		}
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("inspectObjects: encode response: %v", err)
		}

		return

	}

	http.Error(w, "Объект не найден", http.StatusNotFound)
}

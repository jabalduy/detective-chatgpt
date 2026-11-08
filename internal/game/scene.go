package game

import "fmt"

// Scene is a location backdrop. Scene choices are authored as data, not Go handlers.
type Scene struct {
 ID string `json:"id"`
 Title string `json:"title"`
 Background string `json:"background"`
 Description string `json:"description"`
 Choices []Choice `json:"choices"`
}

// Choice represents either a spoken response or a physical action.
type Choice struct {
 ID string `json:"id"`
 Text string `json:"text"`
 Kind string `json:"kind"`
 NextScene string `json:"next_scene,omitempty"`
 Action
}

// SceneView deliberately excludes effects and requirements to avoid spoiling future outcomes.
type SceneView struct {
 ID string `json:"id"`
 Title string `json:"title"`
 Background string `json:"background"`
 Description string `json:"description"`
 Choices []ChoiceView `json:"choices"`
}
type ChoiceView struct {
 ID string `json:"id"`
 Text string `json:"text"`
 Kind string `json:"kind"`
}

func SceneByID(state *GameState, id string) *Scene {
 for i := range state.Scenes {
  if state.Scenes[i].ID == id { return &state.Scenes[i] }
 }
 return nil
}

func CurrentSceneView(state *GameState) (SceneView, error) {
 scene := SceneByID(state, state.CurrentScene)
 if scene == nil { return SceneView{}, fmt.Errorf("scene %q not found", state.CurrentScene) }
 view := SceneView{ID:scene.ID,Title:scene.Title,Background:scene.Background,Description:scene.Description,Choices:make([]ChoiceView,0)}
 for _, choice := range scene.Choices {
  if RequirementsMet(state, choice.Requirements) {
   view.Choices = append(view.Choices,ChoiceView{ID:choice.ID,Text:choice.Text,Kind:choice.Kind})
  }
 }
 return view,nil
}

func ChooseSceneOption(state *GameState, id string) error {
 if state.Ending!=nil {return fmt.Errorf("game finished")}
 if state.ActiveConversation!="" {return fmt.Errorf("finish conversation first")}
 scene := SceneByID(state, state.CurrentScene)
 if scene == nil { return fmt.Errorf("current scene not found") }
 for _, choice := range scene.Choices {
  if choice.ID != id { continue }
  if !CanPerformAction(state,choice.Action) { return fmt.Errorf("choice unavailable") }
  if choice.NextScene != "" && SceneByID(state,choice.NextScene)==nil { return fmt.Errorf("target scene not found") }
  if !PerformAction(state,choice.Action) { return fmt.Errorf("action failed") }
  if choice.NextScene != "" { state.CurrentScene=choice.NextScene }
  return nil
 }
 return fmt.Errorf("choice not found")
}

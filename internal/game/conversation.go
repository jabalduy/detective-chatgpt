package game

import "fmt"

// Conversation is a branching dialogue tree. Nodes can contain spoken or physical choices.
type Conversation struct {
 ID string `json:"id"`
 CharacterID int `json:"character_id"`
 StartNode string `json:"start_node"`
 Nodes []ConversationNode `json:"nodes"`
}
type ConversationNode struct {
 ID string `json:"id"`
 Text string `json:"text"`
 Emotion string `json:"emotion,omitempty"`
 Choices []ConversationChoice `json:"choices"`
}
type ConversationChoice struct {
 ID string `json:"id"`
 Text string `json:"text"`
 Kind string `json:"kind"`
 NextNode string `json:"next_node,omitempty"`
 End bool `json:"end,omitempty"`
 Action
}
type CharacterState struct {
 Location string `json:"location"`
 Emotion string `json:"emotion"`
 Trust int `json:"trust"`
 Present bool `json:"present"`
}
type ConversationView struct {
 ConversationID string `json:"conversation_id"`
 CharacterID int `json:"character_id"`
 NodeID string `json:"node_id"`
 Text string `json:"text"`
 Emotion string `json:"emotion"`
 Choices []ChoiceView `json:"choices"`
}

func conversationByID(s *GameState,id string) *Conversation {
 for i:=range s.Conversations { if s.Conversations[i].ID==id {return &s.Conversations[i]} }
 return nil
}
func nodeByID(c *Conversation,id string) *ConversationNode {
 for i:=range c.Nodes {if c.Nodes[i].ID==id{return &c.Nodes[i]}}
 return nil
}
func StartConversation(s *GameState,id string) error {
 if s.ActiveConversation!="" || s.ActiveDialogue!=nil || s.ActiveEvent!=nil {return fmt.Errorf("another interaction is active")}
 c:=conversationByID(s,id)
 if c==nil {return fmt.Errorf("conversation not found")}
 character,ok:=s.Characters[c.CharacterID]
 if !ok || !character.Present || character.Location!=s.CurrentScene {return fmt.Errorf("character unavailable")}
 if nodeByID(c,c.StartNode)==nil {return fmt.Errorf("start node missing")}
 s.ActiveConversation=id
 s.ConversationNode=c.StartNode
 if s.Clock!=nil {s.Clock.Pause()}
 return nil
}
func CurrentConversationView(s *GameState)(ConversationView,error){
 c:=conversationByID(s,s.ActiveConversation)
 if c==nil {return ConversationView{},fmt.Errorf("no active conversation")}
 node:=nodeByID(c,s.ConversationNode)
 if node==nil {return ConversationView{},fmt.Errorf("node not found")}
 v:=ConversationView{ConversationID:c.ID,CharacterID:c.CharacterID,NodeID:node.ID,Text:node.Text,Emotion:node.Emotion,Choices:[]ChoiceView{}}
 for _,ch:=range node.Choices {if RequirementsMet(s,ch.Requirements){v.Choices=append(v.Choices,ChoiceView{ID:ch.ID,Text:ch.Text,Kind:ch.Kind})}}
 return v,nil
}
func ChooseConversationOption(s *GameState,id string)(bool,error){
 c:=conversationByID(s,s.ActiveConversation)
 if c==nil {return false,fmt.Errorf("no active conversation")}
 node:=nodeByID(c,s.ConversationNode)
 if node==nil {return false,fmt.Errorf("node not found")}
 for _,ch:=range node.Choices {
  if ch.ID!=id {continue}
  if !CanPerformAction(s,ch.Action){return false,fmt.Errorf("choice unavailable")}
  if !ch.End && nodeByID(c,ch.NextNode)==nil {return false,fmt.Errorf("next node missing")}
  if !PerformAction(s,ch.Action){return false,fmt.Errorf("action failed")}
  if ch.End {
   s.ActiveConversation=""
   s.ConversationNode=""
   if s.Clock!=nil {s.Clock.Resume()}
   return true,nil
  }
  s.ConversationNode=ch.NextNode
  return false,nil
 }
 return false,fmt.Errorf("choice not found")
}

package game

import (
 "fmt"
 "sort"
)

// ValidateNarrative checks structural errors that can trap players in a scene
// or dialogue. It does not prescribe a story or require a specific ending.
func ValidateNarrative(scenes []Scene, conversations []Conversation) []error {
 var issues []error
 sceneIDs:=make(map[string]bool,len(scenes))
 for _,scene:=range scenes {
  if scene.ID=="" {issues=append(issues,fmt.Errorf("scene has empty id"));continue}
  if sceneIDs[scene.ID] {issues=append(issues,fmt.Errorf("duplicate scene %q",scene.ID))}
  sceneIDs[scene.ID]=true
 }
 for _,scene:=range scenes {
  choices:=make(map[string]bool)
  for _,choice:=range scene.Choices {
   if choice.ID=="" {issues=append(issues,fmt.Errorf("scene %q has empty choice id",scene.ID))}
   if choices[choice.ID] {issues=append(issues,fmt.Errorf("scene %q has duplicate choice %q",scene.ID,choice.ID))}
   choices[choice.ID]=true
   if choice.NextScene!="" && !sceneIDs[choice.NextScene] {
    issues=append(issues,fmt.Errorf("scene %q choice %q references missing scene %q",scene.ID,choice.ID,choice.NextScene))
   }
  }
 }
 conversationIDs:=make(map[string]bool)
 for _,conversation:=range conversations {
  if conversation.ID=="" {issues=append(issues,fmt.Errorf("conversation has empty id"))}
  if conversationIDs[conversation.ID] {issues=append(issues,fmt.Errorf("duplicate conversation %q",conversation.ID))}
  conversationIDs[conversation.ID]=true
  nodes:=make(map[string]bool)
  for _,node:=range conversation.Nodes {
   if node.ID=="" {issues=append(issues,fmt.Errorf("conversation %q has empty node id",conversation.ID))}
   if nodes[node.ID] {issues=append(issues,fmt.Errorf("conversation %q has duplicate node %q",conversation.ID,node.ID))}
   nodes[node.ID]=true
  }
  if !nodes[conversation.StartNode] {issues=append(issues,fmt.Errorf("conversation %q has missing start node %q",conversation.ID,conversation.StartNode))}
  for _,node:=range conversation.Nodes {
   choiceIDs:=make(map[string]bool)
   for _,choice:=range node.Choices {
    if choice.ID=="" {issues=append(issues,fmt.Errorf("conversation %q node %q has empty choice id",conversation.ID,node.ID))}
    if choiceIDs[choice.ID] {issues=append(issues,fmt.Errorf("conversation %q node %q has duplicate choice %q",conversation.ID,node.ID,choice.ID))}
    choiceIDs[choice.ID]=true
    if !choice.End && !nodes[choice.NextNode] {
     issues=append(issues,fmt.Errorf("conversation %q node %q choice %q references missing node %q",conversation.ID,node.ID,choice.ID,choice.NextNode))
    }
   }
  }
  // Each reachable node must have at least one structural route to an
  // ending. Requirements can still lock choices at runtime.
  reachable:=make(map[string]bool)
  var visit func(string)
  visit=func(id string){
   if reachable[id] || !nodes[id] {return}
   reachable[id]=true
   for _,node:=range conversation.Nodes {if node.ID==id {
    for _,choice:=range node.Choices {if !choice.End {visit(choice.NextNode)}}
    break
   }}
  }
  visit(conversation.StartNode)
  canFinish:=make(map[string]bool)
  for _,node:=range conversation.Nodes {
   for _,choice:=range node.Choices {if choice.End {canFinish[node.ID]=true}}
  }
  changed:=true
  for changed {
   changed=false
   for _,node:=range conversation.Nodes {
    if canFinish[node.ID] {continue}
    for _,choice:=range node.Choices {
     if !choice.End && canFinish[choice.NextNode] {canFinish[node.ID]=true;changed=true;break}
    }
   }
  }
  ids:=make([]string,0,len(reachable))
  for id:=range reachable {if !canFinish[id]{ids=append(ids,id)}}
  sort.Strings(ids)
  for _,id:=range ids {issues=append(issues,fmt.Errorf("conversation %q reachable node %q has no path to an ending",conversation.ID,id))}
 }
 return issues
}

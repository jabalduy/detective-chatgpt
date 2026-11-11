package game

import (
 "strings"
 "testing"
)

func TestValidateEventsAcceptsPlayableUrgentDecision(t *testing.T) {
 events:=[]Event{{ID:1,At:"23:45",TimeLimit:10,DefaultOptionID:2,Options:[]EventOption{{ID:1,Text:"Act"},{ID:2,Text:"Wait"}}}}
 if issues:=ValidateEvents(events);len(issues)!=0 {t.Fatalf("valid event rejected: %v",issues)}
}

func TestValidateEventsReportsUnplayableDefinitions(t *testing.T) {
 events:=[]Event{
  {ID:1,At:"25:99",TimeLimit:10,DefaultOptionID:3,Options:[]EventOption{{ID:1},{ID:1}}},
  {ID:1,TimeLimit:10,DefaultOptionID:2,Options:[]EventOption{{ID:2,Action:Action{Requirements:[]Requirement{{Type:"clue_found",ID:99}}}}}},
  {ID:3,TimeLimit:-1},
 }
 issues:=ValidateEvents(events)
 combined:=""
 for _,issue:=range issues {combined+=issue.Error()+"\n"}
 for _,part:=range []string{"invalid trigger time","duplicate option","no valid default","duplicate event id","default option must not require","negative time limit","no options"} {
  if !strings.Contains(combined,part) {t.Errorf("missing validation %q in %s",part,combined)}
 }
}

package game

import (
 "fmt"
 "time"
)

// ValidateEvents reports event definitions that could strand the operator
// without a valid response. It is intended for case authoring and CI.
func ValidateEvents(events []Event) []error {
 var issues []error
 seenEvents:=make(map[int]bool,len(events))
 for _,e:=range events {
  if e.ID<=0 {issues=append(issues,fmt.Errorf("event id must be positive: %d",e.ID))}
  if seenEvents[e.ID] {issues=append(issues,fmt.Errorf("duplicate event id %d",e.ID))}
  seenEvents[e.ID]=true
  if e.TimeLimit<0 {issues=append(issues,fmt.Errorf("event %d has negative time limit",e.ID))}
  if e.At!="" {
   if _,err:=time.Parse("15:04",e.At);err!=nil {
    issues=append(issues,fmt.Errorf("event %d has invalid trigger time %q",e.ID,e.At))
   }
  }
  options:=make(map[int]bool,len(e.Options))
  for _,o:=range e.Options {
   if o.ID<=0 {issues=append(issues,fmt.Errorf("event %d option id must be positive: %d",e.ID,o.ID))}
   if options[o.ID] {issues=append(issues,fmt.Errorf("event %d has duplicate option %d",e.ID,o.ID))}
   options[o.ID]=true
  }
  if len(e.Options)==0 {issues=append(issues,fmt.Errorf("event %d has no options",e.ID))}
  if e.TimeLimit>0 {
   if e.DefaultOptionID==0 || !options[e.DefaultOptionID] {
    issues=append(issues,fmt.Errorf("timed event %d has no valid default option",e.ID))
   } else {
    for _,o:=range e.Options {
     if o.ID==e.DefaultOptionID && len(o.Requirements)>0 {
      issues=append(issues,fmt.Errorf("timed event %d default option must not require evidence or trust",e.ID))
     }
    }
   }
  }
 }
 return issues
}

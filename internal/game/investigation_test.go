package game

import (
 "testing"
 "time"
)

func TestExaminationCompletesAfterGameTime(t *testing.T){
 start:=time.Date(2026,1,1,22,0,0,0,time.UTC)
 s:=&GameState{Clock:NewClock(start),FoundClues:map[int]bool{7:true},FoundFacts:make(map[int]bool),
 Examinations:[]Examination{{ID:"lab",DurationSeconds:120,Requirements:[]Requirement{{Type:"clue_found",ID:7}},Effects:[]Effect{{Type:"discover_fact",ID:8}},Result:"match"}}}
 if err:=StartExamination(s,"lab");err!=nil{t.Fatal(err)}
 if err:=StartExamination(s,"lab");err==nil{t.Fatal("duplicate start allowed")}
 s.Clock.Add(119*time.Second)
 TickExaminations(s)
 if s.FoundFacts[8]{t.Fatal("completed too soon")}
 s.Clock.Add(time.Second)
 views:=ExaminationViews(s)
 if !s.FoundFacts[8] || len(views)!=1 || views[0].Status!="completed" || views[0].Result!="match" {t.Fatalf("views=%+v",views)}
 if !RequirementsMet(s,[]Requirement{{Type:"examination_completed",Key:"lab"}}){t.Fatal("completion requirement failed")}
}

func TestPresentEvidenceRequiresDiscoveryAndPresence(t *testing.T){
 start:=time.Date(2026,1,1,22,0,0,0,time.UTC)
 s:=&GameState{Clock:NewClock(start),CurrentScene:"corridor",
 Characters:map[int]CharacterState{1:{Location:"corridor",Present:true}},
 FoundFacts:make(map[int]bool),PresentedEvidence:make(map[string]bool),
 EvidencePresentations:[]EvidencePresentation{{ID:"show",CharacterID:1,EvidenceType:"fact",EvidenceID:3,Response:"Interesting",TimeCost:20,Effects:[]Effect{{Type:"character_trust_down",ID:1}}}}}
 if _,err:=PresentEvidence(s,"show");err==nil{t.Fatal("undiscovered evidence accepted")}
 s.FoundFacts[3]=true
 response,err:=PresentEvidence(s,"show")
 if err!=nil || response!="Interesting" {t.Fatalf("response=%s err=%v",response,err)}
 if s.Characters[1].Trust!=-1 || s.Clock.GetTime().Sub(start)!=20*time.Second {t.Fatal("effects or time missing")}
 if _,err:=PresentEvidence(s,"show");err==nil{t.Fatal("duplicate presentation accepted")}
}

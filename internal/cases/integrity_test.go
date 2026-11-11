package cases

import (
 "path/filepath"
 "strings"
 "testing"
 "detective/internal/game"
)

func TestExpress1928Integrity(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 for _,issue:=range ValidateCaseIntegrity(def) {t.Error(issue)}
}

func TestExpress1928IntegrityDetectsBrokenReference(t *testing.T) {
 def,err:=LoadCase(filepath.Join("..","..","cases","case_express_1928"))
 if err!=nil {t.Fatal(err)}
 if len(def.Scenes)==0 || len(def.Scenes[0].Choices)==0 {t.Fatal("missing test scene")}
 def.Scenes[0].Choices[0].Effects=append(def.Scenes[0].Choices[0].Effects,
  game.Effect{Type:"discover_fact",ID:999999})
 found:=false
 for _,issue:=range ValidateCaseIntegrity(def) {
  if strings.Contains(issue.Error(),"unknown fact 999999"){found=true}
 }
 if !found {t.Fatal("integrity validator failed to catch unknown fact")}
}

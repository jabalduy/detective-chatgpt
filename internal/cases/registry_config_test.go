package cases

import (
 "errors"
 "path/filepath"
 "testing"
)

func TestCaseRegistryCustomDirectory(t *testing.T) {
 root:=t.TempDir()
 t.Setenv("GAME_CASE_DIR",root)
 registry:=NewCaseRegistry()
 for _,id:=range []string{"case_express_1928","case_train","case_017","case_dinner"} {
  got,err:=registry.GetPath(id)
  if err!=nil {t.Fatalf("%s: %v",id,err)}
  want:=filepath.Join(root,id)
  if got!=want {t.Errorf("%s path = %q, want %q",id,got,want)}
 }
 if _,err:=registry.GetPath("../case_express_1928");!errors.Is(err,ErrCaseNotFound) {
  t.Fatalf("unregistered case identifier was accepted: %v",err)
 }
}

func TestCaseRegistryDefaultDirectory(t *testing.T) {
 t.Setenv("GAME_CASE_DIR","")
 got,err:=NewCaseRegistry().GetPath("case_express_1928")
 if err!=nil {t.Fatal(err)}
 if got!=filepath.Join("cases","case_express_1928") {
  t.Fatalf("default case directory changed: %q",got)
 }
}

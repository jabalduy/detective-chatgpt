package cases

import (
 "path/filepath"
 "testing"

 "detective/internal/game"
)

func TestExpress1928PrototypeLoads(t *testing.T) {
 path := filepath.Join("..", "..", "cases", "case_express_1928")
 def, err := LoadCase(path)
 if err != nil { t.Fatalf("load prototype: %v", err) }
 if def.ID != "case_express_1928" { t.Fatalf("unexpected case ID %q", def.ID) }
 if def.Solution.SuspectID != 1 || def.Solution.MotiveID != 1 { t.Fatalf("unexpected solution: %+v", def.Solution) }
 if len(def.Suspects) != 8 { t.Fatalf("expected 8 suspects, got %d", len(def.Suspects)) }
 if len(def.Scenes) < 4 { t.Fatalf("expected navigable scenes, got %d", len(def.Scenes)) }
 if issues := game.ValidateNarrative(def.Scenes, def.Conversations); len(issues) > 0 { t.Fatalf("invalid narrative: %v", issues) }

 factIDs := map[int]bool{}
 for _, f := range def.Facts { factIDs[f.ID] = true }
 for _, id := range def.Solution.FactIDs {
  if !factIDs[id] { t.Errorf("solution references missing fact %d", id) }
 }
 for _, d := range def.Deductions {
  for _, id := range d.Required {
   if !factIDs[id] { t.Errorf("deduction %d requires missing fact %d", d.ID, id) }
  }
 }
}

func TestExpress1928Registered(t *testing.T) {
 p, err := NewCaseRegistry().GetPath("case_express_1928")
 if err != nil { t.Fatal(err) }
 if p != "cases/case_express_1928" { t.Fatalf("unexpected path %q", p) }
}

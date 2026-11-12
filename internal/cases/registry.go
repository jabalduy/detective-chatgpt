package cases

import (
	"errors"
	"fmt"
 "os"
 "path/filepath"
)

var ErrCaseNotFound = errors.New("case not found")

type CaseRegistry struct {
	cases map[string]CaseEntry
}

type CaseEntry struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Path  string `json:"-"`
}

func NewCaseRegistry() *CaseRegistry {
	return &CaseRegistry{
		cases: map[string]CaseEntry{
			"case_017": {
				ID:    "case_017",
				Title: "Последний экземпляр",
				Path:  "cases/case_017",
			},

			"case_dinner": {
				ID:    "case_dinner",
				Title: "Ужин на шестерых",
				Path:  "cases/case_dinner",
			},

			"case_express_1928": {
				ID: "case_express_1928",
				Title: "Ночной экспресс: дело Виктора",
				Path: "cases/case_express_1928",
			},

			"case_train": {
				ID:    "case_train",
				Title: "Ночной экспресс",
				Path:  "cases/case_train",
			},
		},
	}
}

func (r *CaseRegistry) GetPath(caseID string) (string, error) {
	entry, ok := r.cases[caseID]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrCaseNotFound, caseID)
	}

	if root:=os.Getenv("GAME_CASE_DIR");root!="" {return filepath.Join(root,filepath.Base(entry.Path)),nil}
 return entry.Path, nil
}

func (r *CaseRegistry) List() []CaseEntry {
	result := make([]CaseEntry, 0, len(r.cases))

	for _, entry := range r.cases {
		result = append(result, entry)
	}

	return result
}

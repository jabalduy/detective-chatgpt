package game

type Clue struct {
	Name  string `json:"name"`
	About string `json:"about"`
	ObjID int    `json:"obj_id"`
	Key   bool   `json:"key"`
}

type Object struct {
	LocID int    `json:"loc_id"`
	Name  string `json:"name"`
	About string `json:"about"`
	ObjID int    `json:"obj_id"`
	Key   bool   `json:"key"`

	InspectAction Action `json:"inspect_action"`

	InitiallyOpen bool `json:"is_open"`
}

type ObjectResponse struct {
	LocID    int    `json:"loc_id"`
	Name     string `json:"name"`
	About    string `json:"about"`
	ObjID    int    `json:"obj_id"`
	Key      bool   `json:"key"`
	Searched bool   `json:"searched"`
}

type Location struct {
	Name   string `json:"name"`
	About  string `json:"about"`
	Object Object `json:"-"`
	LocID  int    `json:"loc_id"`

	InitiallyOpen bool `json:"is_open"`
}

type Suspect struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	About  string `json:"about"`
	SusID  int    `json:"sus_id"`
	IsClue bool   `json:"is_clue"`
}

type Dialogue struct {
	SusID    int    `json:"sus_id"`
	DialID   int    `json:"dial_id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`

	InitiallyOpen bool `json:"is_open"`

	IsClue bool   `json:"is_clue"`
	ObjID  int    `json:"obj_id"`
	Fact   string `json:"fact"`
	Key    bool   `json:"key"`
	Asked  bool   `json:"asked,omitempty"`
	Confronted bool `json:"confronted,omitempty"`

	Confrontation *Confrontation `json:"confrontation,omitempty"`

	Action
}

type DialogueID struct {
	SusID  int
	DialID int
}

func (d Dialogue) ID() DialogueID {
	return DialogueID{
		SusID:  d.SusID,
		DialID: d.DialID,
	}
}

type Confrontation struct {
	Challenge string `json:"challenge"`
	Answer    string `json:"answer"`

	Action
}

type Endings struct {
	Win      string `json:"win"`
	Unsolved string `json:"unsolved"`
	Fail     string `json:"fail"`
}

type CaseInfo struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Intro     string   `json:"intro"`
	KnownInfo string   `json:"known_info"`
	Endings   Endings  `json:"endings"`
	Solution  Solution `json:"solution"`
}

type Solution struct {
	SuspectID int   `json:"suspect_id"`
	MotiveID  int   `json:"motive_id"`
	FactIDs   []int `json:"fact_ids"`
}

type Requirement struct {
	Type   string `json:"type"`
	ID     int    `json:"id"`
	SusID  int    `json:"sus_id,omitempty"`
	DialID int    `json:"dial_id,omitempty"`
}

type Effect struct {
	Type   string `json:"type"`
	ID     int    `json:"id"`
	SusID  int    `json:"sus_id,omitempty"`
	DialID int    `json:"dial_id,omitempty"`
}

type Action struct {
	Requirements []Requirement `json:"requirements"`
	Effects      []Effect      `json:"effects"`
	TimeCost     int           `json:"time_cost"`
}

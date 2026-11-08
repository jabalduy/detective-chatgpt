package game

type Motive struct {
	ID          int    `json:"id"`
	SuspectID   int    `json:"suspect_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

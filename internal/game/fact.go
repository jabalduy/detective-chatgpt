package game

type Fact struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	About string `json:"about"`
	Key   bool   `json:"key"`
}

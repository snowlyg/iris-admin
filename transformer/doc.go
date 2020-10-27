package transformer

type Doc struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	Chapters  []*Chapter `json:"chapters"`
	CreatedAt string     `json:"created_at"`
}

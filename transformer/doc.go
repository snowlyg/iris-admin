package transformer

type Doc struct {
	Id         int        `json:"id"`
	ChapterMun int64      `json:"chapter_mun"`
	Name       string     `json:"name"`
	Chapters   []*Chapter `json:"chapters"`
	CreatedAt  string     `json:"created_at"`
}

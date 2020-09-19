package transformer

type Article struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	ImageUri     string `json:"image_uri"`
	SourceUri    string `json:"source_uri"`
	IsOriginal   bool   `json:"is_original"`
	Content      string `json:"content"`
	ContentShort string `json:"content_short"`
	Status       string `json:"status"`
	DisplayTime  string `json:"display_time"`
	CreatedAt    string `json:"created_at"`
}

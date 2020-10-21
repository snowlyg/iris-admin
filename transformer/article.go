package transformer

type Article struct {
	Id           int      `json:"id"`
	Title        string   `json:"title"`
	Author       string   `json:"author"`
	ImageUri     string   `json:"image_uri"`
	SourceUri    string   `json:"source_uri"`
	IsOriginal   bool     `json:"is_original"`
	Content      string   `json:"content"`
	ContentShort string   `json:"content_short"`
	Status       string   `json:"status"`
	TagNames     []string `json:"tag_names"`
	Type         *Type    `json:"type"`
	DisplayTime  string   `json:"display_time"`
	DisplayAt    string   `time:"2006-01-02 15:04:05" name:"DisplayTime" json:"display_at"`
	CreatedAt    string   `json:"created_at"`
}

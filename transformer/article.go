package transformer

type Article struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Cover      string `json:"cover"`
	Source     string `json:"source"`
	IsOriginal string `json:"is_original"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

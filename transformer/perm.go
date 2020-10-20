package transformer

type Permission struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Act         string `json:"act"`
	CreatedAt   string `json:"created_at"`
}

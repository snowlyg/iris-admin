package transformer

type Permission struct {
	Id          int    `json:"id"`
	Name        string `json:"key"`
	DisplayName string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

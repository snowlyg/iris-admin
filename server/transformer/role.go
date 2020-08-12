package transformer

type Role struct {
	Id          int           `json:"id"`
	Name        string        `json:"key"`
	DisplayName string        `json:"name"`
	Description string        `json:"description"`
	Perms       []*Permission `json:"routes"`
	CreatedAt   string        `json:"created_at"`
}

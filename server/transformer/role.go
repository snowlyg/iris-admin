package transformer

type Role struct {
	Id          int           `json:"id"`
	Name        string        `json:"name"`
	DisplayName string        `json:"display_name"`
	Description string        `json:"description"`
	Perms       []*Permission `json:"perms"`
	CreatedAt   string        `json:"created_at"`
}

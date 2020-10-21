package transformer

type User struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Username  string  `json:"username"`
	Intro     string  `json:"introduction"`
	Avatar    string  `json:"avatar"`
	Roles     []*Role `json:"roles"`
	RoleIds   []int   `json:"role_ids"`
	CreatedAt string  `json:"created_at"`
}

package transformer

type User struct {
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	Username  string   `json:"username"`
	Intro     string   `json:"introduction"`
	Avatar    string   `json:"avatar"`
	RoleName  []string `json:"roles"`
	RoleIds   []int    `json:"role_ids"`
	CreatedAt string   `json:"created_at"`
}

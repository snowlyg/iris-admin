package transformer

type User struct {
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	Username  string   `json:"username"`
	RoleName  []string `json:"roles"`
	RoleIds   []int    `json:"role_ids"`
	CreatedAt string   `json:"created_at"`
}

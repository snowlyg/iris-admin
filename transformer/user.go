package transformer

type User struct {
	Id        int
	Name      string
	Username  string
	RoleName  string `gtf:"Role.DisplayName"`
	CreatedAt string
}

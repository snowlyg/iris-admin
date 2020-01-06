package transformer

type User struct {
	Name      string
	Username  string
	RoleName  string `gtf:"Role.DisplayName"`
	CreatedAt string
}

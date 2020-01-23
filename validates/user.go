package validates

type CreateUpdateUserRequest struct {
	Username string `json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required"  comment:"密码"`
	Name     string `json:"name" validate:"required,gte=2,lte=50"  comment:"名称"`
	RoleIds  []uint `json:"role_ids"  validate:"required" comment:"角色"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required"  comment:"密码"`
}

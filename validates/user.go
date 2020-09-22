package validates

type LoginRequest struct {
	Username string `json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required"  comment:"密码"`
}

package models

type AdminUserLogin struct {
	Username string `json:"username" validate:"required,gte=4,lte=50"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required,gte=4,lte=50"`
	Phone    string `json:"phone" validate:"required"`
	RoleId   uint   `json:"role_id" validate:"required"`
}

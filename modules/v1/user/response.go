package user

import (
	"github.com/snowlyg/iris-admin/g"
)

type Response struct {
	g.Model
	BaseUser
	Roles []string `gorm:"-" json:"roles"`
}

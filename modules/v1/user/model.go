package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	BaseUser
	RoleIds []uint `gorm:"-" json:"role_ids"`
}

type BaseUser struct {
	Name     string `gorm:"index;not null; type:varchar(60)" json:"name" `
	Username string `gorm:"uniqueIndex;not null;type:varchar(60)" json:"username"`
	Password string `gorm:"type:varchar(250)" json:"password"`
	Intro    string `gorm:"not null; type:varchar(512)" json:"introduction"`
	Avatar   string `gorm:"type:varchar(1024)" json:"avatar"`
}

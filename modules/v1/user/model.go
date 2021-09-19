package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	BaseUser
	Password string `gorm:"type:varchar(250)" json:"password"`
	RoleIds  []uint `gorm:"-" json:"role_ids"`
}

type BaseUser struct {
	Name     string `gorm:"index;not null; type:varchar(60)" json:"name" `
	Username string `gorm:"uniqueIndex;not null;type:varchar(60)" json:"username"`
	Intro    string `gorm:"not null; type:varchar(512)" json:"intro"`
	Avatar   string `gorm:"type:varchar(1024)" json:"avatar"`
}

type Avatar struct {
	Avatar string `json:"avatar"`
}

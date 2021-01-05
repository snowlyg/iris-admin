package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name     string `gorm:"index;not null; type:varchar(60)" json:"name" `
	Username string `gorm:"uniqueIndex;not null;type:varchar(60)" json:"username"`
	Password string `gorm:"type:varchar(100)" json:"password"`
	Intro    string `gorm:"not null; type:varchar(512)" json:"introduction"`
	Avatar   string `gorm:"type:longText" json:"avatar"`

	Oplogs  []Oplog
	RoleIds []uint `gorm:"-" json:"role_ids"`
}

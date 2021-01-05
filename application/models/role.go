package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model

	Name        string     `gorm:"uniqueIndex;not null; type:varchar(256)" json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	DisplayName string     `gorm:"type:varchar(256)" json:"display_name" comment:"显示名称"`
	Description string     `gorm:"type:varchar(256)" json:"description" comment:"描述"`
	Perms       [][]string `gorm:"-" json:"perms" comment:"权限 name, act "`
}

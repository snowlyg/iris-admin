package models

import (
	"gorm.io/gorm"
)

type Oplog struct {
	gorm.Model
	ModelName  string `gorm:"type:varchar(256)" json:"model_name" validate:"required" comment:"模块名称"`
	ActionName string `gorm:"type:varchar(256)" json:"action_name" validate:"required" comment:"操作名称"`
	Content    string `gorm:"type:longText" json:"content"  comment:"操作内容"`
	UserID     uint
}

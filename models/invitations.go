package models

import (
	"github.com/jinzhu/gorm"
)

type Invitations struct {
	gorm.Model
	Code      string `gorm:"not null comment('分享授权码') VARCHAR(191)"`
	DepartIds string `gorm:"not null comment('授权部门id数组') VARCHAR(191)"`
	ReportId  int    `gorm:"not null index INT(10)"`
}

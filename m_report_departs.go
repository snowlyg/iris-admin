package main

import (
	"github.com/jinzhu/gorm"
)

type ReportDeparts struct {
	gorm.Model
	Name     string `gorm:"not null comment('库名部门名称') VARCHAR(191)"`
	Icon     string `gorm:"comment('库名部门图标') VARCHAR(191)"`
	ReportId int    `gorm:"not null index INT(10)"`
}

package models

import (
	"IrisYouQiKangApi/system"
	"github.com/jinzhu/gorm"
)

type PlanDeparts struct {
	gorm.Model
	Name   string `gorm:"not null comment('库名部门名称') VARCHAR(191)"`
	Icon   string `gorm:"comment('库名部门图标') VARCHAR(191)"`
	PlanId int    `gorm:"not null index INT(10)"`
}

func init() {
	system.DB.AutoMigrate(&PlanDeparts{})
}

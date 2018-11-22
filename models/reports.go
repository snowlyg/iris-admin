package models

import "github.com/jinzhu/gorm"

type Reports struct {
	gorm.Model
	Name    string `gorm:"not null comment('报告名称') VARCHAR(191)"`
	OrderId int    `gorm:"not null index INT(10)"`
}

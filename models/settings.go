package models

import (
	"IrisYouQiKangApi/system"
	"github.com/jinzhu/gorm"
)

type Settings struct {
	gorm.Model
	Key   string `gorm:"not null unique VARCHAR(191)"`
	Value string `gorm:"not null VARCHAR(191)"`
}

func init() {
	system.DB.AutoMigrate(&Settings{})
}

package models

import (
	"IrisYouQiKangApi/system"
	"github.com/jinzhu/gorm"
)

type Messages struct {
	gorm.Model
	Content string `gorm:"not null comment('微信消息内容') LONGTEXT"`
}

func init() {
	system.DB.AutoMigrate(&Messages{})
}

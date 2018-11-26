package models

import (
	"IrisYouQiKangApi/system"
	"github.com/jinzhu/gorm"
)

type Revisions struct {
	gorm.Model
	RevisionableType string `gorm:"not null index(revisions_revisionable_id_revisionable_type_index) VARCHAR(191)"`
	RevisionableId   int    `gorm:"not null index(revisions_revisionable_id_revisionable_type_index) INT(11)"`
	UserId           int    `gorm:"INT(11)"`
	Key              string `gorm:"not null VARCHAR(191)"`
	OldValue         string `gorm:"LONGTEXT"`
	NewValue         string `gorm:"LONGTEXT"`
	Device           string `gorm:"VARCHAR(191)"`
	Ip               string `gorm:"VARCHAR(191)"`
	DeviceType       string `gorm:"VARCHAR(191)"`
	Address          string `gorm:"VARCHAR(191)"`
	Browser          string `gorm:"VARCHAR(191)"`
	Platform         string `gorm:"VARCHAR(191)"`
	Language         string `gorm:"VARCHAR(191)"`
}

func init() {
	system.DB.AutoMigrate(&Revisions{})
}

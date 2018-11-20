package models

import (
	"time"
)

type Revisions struct {
	Id               int       `xorm:"not null pk autoincr INT(10)"`
	RevisionableType string    `xorm:"not null index(revisions_revisionable_id_revisionable_type_index) VARCHAR(191)"`
	RevisionableId   int       `xorm:"not null index(revisions_revisionable_id_revisionable_type_index) INT(11)"`
	UserId           int       `xorm:"INT(11)"`
	Key              string    `xorm:"not null VARCHAR(191)"`
	OldValue         string    `xorm:"LONGTEXT"`
	NewValue         string    `xorm:"LONGTEXT"`
	Device           string    `xorm:"VARCHAR(191)"`
	Ip               string    `xorm:"VARCHAR(191)"`
	DeviceType       string    `xorm:"VARCHAR(191)"`
	Address          string    `xorm:"VARCHAR(191)"`
	Browser          string    `xorm:"VARCHAR(191)"`
	Platform         string    `xorm:"VARCHAR(191)"`
	Language         string    `xorm:"VARCHAR(191)"`
	CreatedAt        time.Time `xorm:"TIMESTAMP"`
	UpdatedAt        time.Time `xorm:"TIMESTAMP"`
}

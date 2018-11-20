package models

import (
	"time"
)

type Invitations struct {
	Id        int       `xorm:"not null pk autoincr INT(10)"`
	Code      string    `xorm:"not null comment('分享授权码') VARCHAR(191)"`
	DepartIds string    `xorm:"not null comment('授权部门id数组') VARCHAR(191)"`
	ReportId  int       `xorm:"not null index INT(10)"`
	CreatedAt time.Time `xorm:"TIMESTAMP"`
	UpdatedAt time.Time `xorm:"TIMESTAMP"`
	DeletedAt time.Time `xorm:"TIMESTAMP"`
}

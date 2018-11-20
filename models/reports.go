package models

import (
	"time"
)

type Reports struct {
	Id        int       `xorm:"not null pk autoincr INT(10)"`
	Name      string    `xorm:"not null comment('报告名称') VARCHAR(191)"`
	OrderId   int       `xorm:"not null index INT(10)"`
	CreatedAt time.Time `xorm:"TIMESTAMP"`
	UpdatedAt time.Time `xorm:"TIMESTAMP"`
}

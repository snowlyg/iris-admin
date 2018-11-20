package models

import (
	"time"
)

type PlanDeparts struct {
	Id        int       `xorm:"not null pk autoincr INT(10)"`
	Name      string    `xorm:"not null comment('库名部门名称') VARCHAR(191)"`
	Icon      string    `xorm:"comment('库名部门图标') VARCHAR(191)"`
	PlanId    int       `xorm:"not null index INT(10)"`
	CreatedAt time.Time `xorm:"TIMESTAMP"`
	UpdatedAt time.Time `xorm:"TIMESTAMP"`
	DeletedAt time.Time `xorm:"TIMESTAMP"`
}

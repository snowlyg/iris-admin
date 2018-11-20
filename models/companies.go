package models

import (
	"time"
)

type Companies struct {
	Id        int       `xorm:"not null pk autoincr INT(10)"`
	Name      string    `xorm:"not null comment('客户名称') VARCHAR(191)"`
	Creator   string    `xorm:"not null comment('创建人') VARCHAR(191)"`
	Logo      string    `xorm:"comment('logo') VARCHAR(191)"`
	Preview   int       `xorm:"default 0 comment('设置演示数据') INT(1)"`
	CreatedAt time.Time `xorm:"TIMESTAMP"`
	UpdatedAt time.Time `xorm:"TIMESTAMP"`
	DeletedAt time.Time `xorm:"TIMESTAMP"`
}

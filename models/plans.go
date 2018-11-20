package models

import (
	"time"
)

type Plans struct {
	Id        int       `xorm:"not null pk autoincr INT(10)"`
	Name      string    `xorm:"not null comment('库名称') VARCHAR(191)"`
	Editer    string    `xorm:"not null comment('编辑人') VARCHAR(191)"`
	IsParent  int       `xorm:"not null default 0 comment('是否是标准库') TINYINT(1)"`
	DeletedAt time.Time `xorm:"TIMESTAMP"`
	CreatedAt time.Time `xorm:"TIMESTAMP"`
	UpdatedAt time.Time `xorm:"TIMESTAMP"`
}

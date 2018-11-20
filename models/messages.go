package models

import (
	"time"
)

type Messages struct {
	Id        int       `xorm:"not null pk autoincr INT(10)"`
	Content   string    `xorm:"not null comment('微信消息内容') LONGTEXT"`
	CreatedAt time.Time `xorm:"TIMESTAMP"`
	UpdatedAt time.Time `xorm:"TIMESTAMP"`
}

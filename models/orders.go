package models

import (
	"time"
)

type Orders struct {
	Id          int       `xorm:"not null pk autoincr INT(10)"`
	Name        string    `xorm:"not null comment('订单名称') VARCHAR(191)"`
	Status      string    `xorm:"not null comment('订单状态（‘未开始’，‘进行中’，‘已完成’）') VARCHAR(191)"`
	StartAt     time.Time `xorm:"comment('启动时间') DATETIME"`
	OrderNumber string    `xorm:"not null comment('订单号') VARCHAR(191)"`
	PlanId      int       `xorm:"not null index INT(10)"`
	CompanyId   int       `xorm:"not null index INT(10)"`
	CreatedAt   time.Time `xorm:"TIMESTAMP"`
	UpdatedAt   time.Time `xorm:"TIMESTAMP"`
	DeletedAt   time.Time `xorm:"TIMESTAMP"`
}

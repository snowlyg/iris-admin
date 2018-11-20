package models

import (
	"time"
)

type Permissions struct {
	Id          int       `xorm:"not null pk autoincr INT(10)"`
	Name        string    `xorm:"not null VARCHAR(191)"`
	GuardName   string    `xorm:"not null VARCHAR(191)"`
	DisplayName string    `xorm:"VARCHAR(191)"`
	Description string    `xorm:"VARCHAR(191)"`
	DeletedAt   time.Time `xorm:"TIMESTAMP"`
	CreatedAt   time.Time `xorm:"TIMESTAMP"`
	UpdatedAt   time.Time `xorm:"TIMESTAMP"`
}

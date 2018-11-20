package models

import (
	"time"
)

type OauthToken struct {
	Id        int       `xorm:"not null pk autoincr INT(11)"`
	Token     string    `xorm:"not null default '' comment('Token') VARCHAR(191)"`
	AppId     string    `xorm:"not null default '' comment('Appid') VARCHAR(191)"`
	Secret    string    `xorm:"not null default '' comment('Secret') VARCHAR(191)"`
	ExpressIn int64     `xorm:"not null default 0 comment('是否是标准库') BIGINT(20)"`
	CreatedAt time.Time `xorm:"DATETIME"`
	UpdatedAt time.Time `xorm:"DATETIME"`
}

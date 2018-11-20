package models

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Users struct {
	Id               int       `xorm:"not null pk autoincr INT(10)"`
	Name             string    `xorm:"not null VARCHAR(191)"`
	Username         string    `xorm:"VARCHAR(191)"`
	Password         string    `xorm:"not null VARCHAR(191)"`
	Confirmed        int       `xorm:"not null default 0 TINYINT(1)"`
	IsClient         int       `xorm:"not null default 0 TINYINT(1)"`
	IsFrozen         int       `xorm:"not null default 0 TINYINT(1)"`
	IsAudit          int       `xorm:"not null default 0 TINYINT(1)"`
	IsClientAdmin    int       `xorm:"not null default 0 TINYINT(1)"`
	WechatName       string    `xorm:"VARCHAR(191)"`
	WechatAvatar     string    `xorm:"VARCHAR(191)"`
	Email            string    `xorm:"unique VARCHAR(191)"`
	OpenId           string    `xorm:"unique VARCHAR(191)"`
	WechatVerfiyTime time.Time `xorm:"DATETIME"`
	IsWechatVerfiy   int       `xorm:"not null default 0 TINYINT(1)"`
	Phone            string    `xorm:"unique VARCHAR(191)"`
	RoleId           int       `xorm:"INT(10)"`
	RememberToken    string    `xorm:"VARCHAR(100)"`
	CreatedAt        time.Time `xorm:"TIMESTAMP"`
	UpdatedAt        time.Time `xorm:"TIMESTAMP"`
	DeletedAt        time.Time `xorm:"TIMESTAMP"`
}

func UserAdminCheckLogin(username, password string) (u *Users, err error) {

	u = &Users{Username: username}
	_, err = DB.Get(u)

	if err != nil {
		return
	}

	//hashPassword,err := bcrypt.GenerateFromPassword()
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	if err != nil {
		return
	}

	return
}

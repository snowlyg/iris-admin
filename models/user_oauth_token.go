package models

type UserOauthToken struct {
	Users      `xorm:"extends"`
	OauthToken `xorm:"extends"`
}

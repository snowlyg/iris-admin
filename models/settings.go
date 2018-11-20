package models

type Settings struct {
	Id    int    `xorm:"not null pk autoincr INT(10)"`
	Key   string `xorm:"not null unique VARCHAR(191)"`
	Value string `xorm:"not null VARCHAR(191)"`
}

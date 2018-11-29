package main

import (
	"github.com/jinzhu/gorm"
)

type Roles struct {
	gorm.Model
	Name        string `gorm:"not null VARCHAR(191)"`
	GuardName   string `gorm:"not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Description string `gorm:"VARCHAR(191)"`
	Level       int    `gorm:"not null default 0 INT(10)"`
}

/**
 * 获取所有的账号
 * @method GetAllRoles
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func MGetAllRoles(kw string, cp int, mp int) (roles []*Roles) {
	if len(kw) > 0 {
		db.Model(Roles{}).Where("name=?", kw).Offset(cp - 1).Limit(mp).Find(&roles)
	}
	db.Model(Roles{}).Offset(cp - 1).Limit(mp).Find(&roles)

	return
}

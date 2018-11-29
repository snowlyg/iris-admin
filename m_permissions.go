package main

import (
	"github.com/jinzhu/gorm"
)

type Permissions struct {
	gorm.Model
	Name        string `gorm:"not null VARCHAR(191)"`
	GuardName   string `gorm:"not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Description string `gorm:"VARCHAR(191)"`
}

/**
 * 获取所有的账号
 * @method GetAllPerms
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func MGetAllPerms(kw string, cp int, mp int) (perms []*Permissions) {

	if len(kw) > 0 {
		db.Model(Permissions{}).Where("name=?", kw).Offset(cp - 1).Limit(mp).Find(&perms)
	}
	db.Model(Permissions{}).Offset(cp - 1).Limit(mp).Find(&perms)

	return
}

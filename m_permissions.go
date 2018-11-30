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
 * @method MGetAllPerms
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func MGetAllPerms(name, orderBy string, offset, limit int) (perms []*Permissions) {
	searchKeys := make(map[string]interface{})
	searchKeys["name"] = name

	MGetAll(searchKeys, orderBy, "", offset, limit).Find(&perms)
	return
}

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
 * @method MGetAllRoles
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func MGetAllRoles(name, orderBy string, offset, limit int) (roles []*Roles) {
	searchKeys := make(map[string]interface{})
	searchKeys["name"] = name
	searchKeys["is_client"] = false

	MGetAll(searchKeys, orderBy, "", offset, limit).Find(&roles)
	return
}

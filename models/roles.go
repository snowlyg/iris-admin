package models

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
func GetAllRoles(kw string, cp int, mp int) (aj ApiJson) {
	roles := make([]Roles, 0)
	if len(kw) > 0 {
		DB.Model(Roles{}).Where("name=?", kw).Offset(cp - 1).Limit(mp).Find(&roles)
	}
	DB.Model(Roles{}).Offset(cp - 1).Limit(mp).Find(&roles)

	auts := TransFormRoles(roles)

	aj.State = true
	aj.Data = auts
	aj.Msg = "操作成功"

	return
}

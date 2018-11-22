package models

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
func GetAllPerms(kw string, cp int, mp int) (aj ApiJson) {
	perms := make([]Permissions, 0)
	if len(kw) > 0 {
		DB.Model(Permissions{}).Where("name=?", kw).Offset(cp - 1).Limit(mp).Find(&perms)
	}
	DB.Model(Permissions{}).Offset(cp - 1).Limit(mp).Find(&perms)

	auts := TransFormPerms(perms)

	aj.State = true
	aj.Data = auts
	aj.Msg = "操作成功"

	return
}

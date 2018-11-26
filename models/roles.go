package models

import (
	"IrisYouQiKangApi/system"
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

func init() {
	system.DB.AutoMigrate(&Roles{})
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
		system.DB.Model(Roles{}).Where("name=?", kw).Offset(cp - 1).Limit(mp).Find(&roles)
	}
	system.DB.Model(Roles{}).Offset(cp - 1).Limit(mp).Find(&roles)

	auts := TransFormRoles(roles)

	aj.Status = true
	aj.Data = auts
	aj.Msg = "操作成功"

	return
}

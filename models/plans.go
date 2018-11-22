package models

import (
	"github.com/jinzhu/gorm"
)

type Plans struct {
	gorm.Model

	Name     string `gorm:"not null comment('库名称') VARCHAR(191)"`
	Editer   string `gorm:"not null comment('编辑人') VARCHAR(191)"`
	IsParent int    `gorm:"not null default 0 comment('是否是标准库') TINYINT(1)"`
}

/**
 * 获取所有的诊断方案
 * @method GetAllPlans
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func GetAllPlans(kw string, cp int, mp int) (aj ApiJson) {
	orders := make([]Plans, 0)
	if len(kw) > 0 {
		DB.Model(Plans{}).Where("is_parent=?", 0).Where("name=?", kw).Offset(cp - 1).Limit(mp).Find(&orders)
	}
	DB.Model(Plans{}).Where("is_parent=?", 0).Offset(cp - 1).Limit(mp).Find(&orders)

	auts := TransFormPlans(orders)

	aj.State = true
	aj.Data = auts
	aj.Msg = "操作成功"

	return
}

/**
 * 获取所有的诊断方案
 * @method GetAllPlans
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func GetAllParentPlans(kw string, cp int, mp int) (aj ApiJson) {
	orders := make([]Plans, 0)
	if len(kw) > 0 {
		DB.Model(Plans{}).Where("is_parent=?", 1).Where("name=?", kw).Offset(cp - 1).Limit(mp).First(&orders)
	}
	DB.Model(Plans{}).Where("is_parent=?", 1).Offset(cp - 1).Limit(mp).Find(&orders)

	auts := TransFormPlans(orders)

	aj.State = true
	aj.Data = auts
	aj.Msg = "操作成功"

	return
}

package main

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
 * @method MGetAllPlans
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func MGetAllPlans(name, orderBy string, offset, limit int) (plans []*Plans) {
	searchKeys := make(map[string]interface{})
	searchKeys["name"] = name
	searchKeys["is_parent"] = false

	MGetAll(searchKeys, orderBy, "", offset, limit).Find(&plans)
	return
}

/**
 * 获取所有的诊断方案
 * @method MGetAllParentPlans
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func MGetAllParentPlans(name, orderBy string, offset, limit int) (plans []*Plans) {
	searchKeys := make(map[string]interface{})
	searchKeys["name"] = name
	searchKeys["is_parent"] = true

	MGetAll(searchKeys, orderBy, "", offset, limit).Find(&plans)
	return
}

package main

import (
	"github.com/jinzhu/gorm"
)

type Companies struct {
	gorm.Model
	Name    string `gorm:"not null comment('客户名称') VARCHAR(191)"`
	Creator string `gorm:"not null comment('创建人') VARCHAR(191)"`
	Logo    string `gorm:"comment('logo') VARCHAR(191)"`
	Preview int    `gorm:"default 0 comment('设置演示数据') INT(1)"`
}

/**
 * 获取所有的客户
 * @method GetAllCompanies
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func MGetAllCompanies(name, orderBy string, offset, limit int) (companies []*Companies) {
	searchKeys := make(map[string]interface{})
	searchKeys["name"] = name

	MGetAll(searchKeys, orderBy, "", offset, limit).Find(&companies)
	return
}

/**
 * 获取所有的客户数量
 * @method GetCompanyCounts
 * @return  {[type]} count int    [description]
 */
func MGetCompanyCounts() (counts int) {
	db.Model(&Companies{}).Count(&counts)
	return
}

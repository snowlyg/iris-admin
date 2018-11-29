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
func MGetAllCompanies(kw string, cp int, mp int) (companies []*Companies) {
	if len(kw) > 0 {
		db.Model(Companies{}).Where("name=?", kw).Offset(cp - 1).Limit(mp).Find(&companies)
	}
	db.Model(Companies{}).Offset(cp - 1).Limit(mp).Find(&companies)

	return
}

/**
 * 获取所有的客户数量
 * @method GetCompanyCounts
 * @return  {[type]} count int    [description]
 */
func MGetCompanyCounts() (count int) {
	db.Model(&Companies{}).Count(&count)
	return
}

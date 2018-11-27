package models

import (
	"IrisYouQiKangApi/system"
	"github.com/jinzhu/gorm"
)

type Companies struct {
	gorm.Model
	Name    string `gorm:"not null comment('客户名称') VARCHAR(191)"`
	Creator string `gorm:"not null comment('创建人') VARCHAR(191)"`
	Logo    string `gorm:"comment('logo') VARCHAR(191)"`
	Preview int    `gorm:"default 0 comment('设置演示数据') INT(1)"`
}

func init() {
	system.DB.AutoMigrate(&Companies{})
}

/**
 * 获取所有的客户
 * @method GetAllCompanies
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func GetAllCompanies(kw string, cp int, mp int) (companies []*Companies) {
	if len(kw) > 0 {
		system.DB.Model(Companies{}).Where("name=?", kw).Offset(cp - 1).Limit(mp).Find(&companies)
	}
	system.DB.Model(Companies{}).Offset(cp - 1).Limit(mp).Find(&companies)

	return
}

package models

import "github.com/jinzhu/gorm"

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
func GetAllCompanies(kw string, cp int, mp int) (aj ApiJson) {
	companies := make([]Companies, 0)
	if len(kw) > 0 {
		DB.Model(Companies{}).Where("name=?", kw).Offset(cp - 1).Limit(mp).Find(&companies)
	}
	DB.Model(Companies{}).Offset(cp - 1).Limit(mp).Find(&companies)

	auts := TransFormCompanies(companies)

	aj.State = true
	aj.Data = auts
	aj.Msg = "操作成功"

	return
}

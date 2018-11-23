package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Orders struct {
	gorm.Model
	Name        string     `gorm:"not null comment('订单名称') VARCHAR(191)"`
	Status      string     `gorm:"not null comment('订单状态（‘未开始’，‘进行中’，‘已完成’）') VARCHAR(191)"`
	StartAt     *time.Time `gorm:"comment('启动时间') DATETIME"`
	OrderNumber string     `gorm:"not null comment('订单号') VARCHAR(191)"`
	PlanId      int        `gorm:"not null index INT(10)"`
	CompanyId   int        `gorm:"not null index INT(10)"`
}

/**
 * 获取所有的订单
 * @method GetAllOrders
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func GetAllOrders(kw string, cp int, mp int) (aj ApiJson) {
	orders := make([]Orders, 0)
	if len(kw) > 0 {
		DB.Model(Orders{}).Where("name=?", kw).Offset(cp - 1).Limit(mp).Find(&orders)
	}
	DB.Model(Orders{}).Offset(cp - 1).Limit(mp).Find(&orders)

	auts := TransFormOrders(orders)

	aj.Status = true
	aj.Data = auts
	aj.Msg = "操作成功"

	return
}

package main

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
func MGetAllOrders(kw string, cp int, mp int) (orders []*Orders) {

	if len(kw) > 0 {
		db.Model(Orders{}).Where("name=?", kw).Offset(cp - 1).Limit(mp).Find(&orders)
	}
	db.Model(Orders{}).Offset(cp - 1).Limit(mp).Find(&orders)

	return
}

/**
 * 获取所有的订单数量
 * @method GetOrderCounts
 * @return  {[type]} counts int    [description]
 */
func MGetOrderCounts() (counts int) {
	db.Model(&Orders{}).Count(&counts)

	return
}

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
 * @method MGetAllOrders
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func MGetAllOrders(name, orderBy string, offset, limit int) (orders []*Orders) {
	searchKeys := make(map[string]interface{})
	searchKeys["name"] = name

	MGetAll(searchKeys, orderBy, "", offset, limit).Find(&orders)
	return
}

/**
 * 获取所有的订单数量
 * @method MGetOrderCounts
 * @return  {[type]} counts int    [description]
 */
func MGetOrderCounts() (counts int) {
	db.Model(&Orders{}).Count(&counts)
	return
}

package controllers

import (
	"IrisYouQiKangApi/models"
	"github.com/kataras/iris"
)

/**
* @api {get} /orders 获取所有的订单
* @apiName 获取所有的订单
* @apiGroup Orders
* @apiVersion 1.0.0
* @apiDescription 获取所有的订单
* @apiSampleRequest /orders
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func GetAllOrders(ctx iris.Context) {
	cp := Tools.ParseInt(ctx.FormValue("cp"), 1)
	mp := Tools.ParseInt(ctx.FormValue("mp"), 20)
	kw := ctx.FormValue("kw")

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(models.GetAllOrders(kw, cp, mp))
}

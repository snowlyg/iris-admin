package main

import (
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
func CGetAllOrders(ctx iris.Context) {
	offset := t.ParseInt(ctx.FormValue("offset"), 1)
	limit := t.ParseInt(ctx.FormValue("limit"), 20)
	name := ctx.FormValue("name")
	orderBy := ctx.FormValue("orderBy")
	orders := MGetAllOrders(name, orderBy, offset, limit)

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(apiResource(true, orders, "操作成功"))
}

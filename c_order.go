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

	cp := t.ParseInt(ctx.FormValue("cp"), 1)
	mp := t.ParseInt(ctx.FormValue("mp"), 20)
	kw := ctx.FormValue("kw")
	orders := MGetAllOrders(kw, cp, mp)

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(apiResource(true, orders, "操作成功"))
}

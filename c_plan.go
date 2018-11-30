package main

import (
	"github.com/kataras/iris"
)

/**
* @api {get} /plans/parent 获取所有的诊断方案
* @apiName 获取所有的诊断方案
* @apiGroup Plans
* @apiVersion 1.0.0
* @apiDescription 获取所有的诊断方案
* @apiSampleRequest /plans/parent
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func CGetAllParentPlans(ctx iris.Context) {
	offset := t.ParseInt(ctx.FormValue("offset"), 1)
	limit := t.ParseInt(ctx.FormValue("limit"), 20)
	name := ctx.FormValue("name")
	orderBy := ctx.FormValue("orderBy")

	plans := MGetAllParentPlans(name, orderBy, offset, limit)

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(apiResource(true, plans, "操作成功"))
}

/**
* @api {get} /plans 获取所有的诊断方案
* @apiName 获取所有的诊断方案
* @apiGroup Plans
* @apiVersion 1.0.0
* @apiDescription 获取所有的诊断方案
* @apiSampleRequest /plans
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func CGetAllPlans(ctx iris.Context) {
	offset := t.ParseInt(ctx.FormValue("offset"), 1)
	limit := t.ParseInt(ctx.FormValue("limit"), 20)
	name := ctx.FormValue("name")
	orderBy := ctx.FormValue("orderBy")

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(MGetAllPlans(name, orderBy, offset, limit))
}

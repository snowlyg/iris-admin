package main

import (
	"github.com/kataras/iris"
)

/**
* @api {get} /companies 获取所有的客户（公司信息）
* @apiName 获取所有的客户（公司信息）
* @apiGroup Companies
* @apiVersion 1.0.0
* @apiDescription 获取所有的客户（公司信息）
* @apiSampleRequest /companies
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func CGetAllCompanies(ctx iris.Context) {
	offset := t.ParseInt(ctx.FormValue("offset"), 1)
	limit := t.ParseInt(ctx.FormValue("limit"), 20)
	name := ctx.FormValue("name")
	orderBy := ctx.FormValue("orderBy")

	ctx.StatusCode(iris.StatusOK)
	companies := MGetAllCompanies(name, orderBy, offset, limit)
	ctx.JSON(apiResource(true, companies, "操作成功"))
}

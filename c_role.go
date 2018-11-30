package main

import (
	"github.com/kataras/iris"
)

/**
* @api {get} /roles 获取所有的角色
* @apiName 获取所有的角色
* @apiGroup Roles
* @apiVersion 1.0.0
* @apiDescription 获取所有的角色
* @apiSampleRequest /roles
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func CGetAllRoles(ctx iris.Context) {
	offset := t.ParseInt(ctx.FormValue("offset"), 1)
	limit := t.ParseInt(ctx.FormValue("limit"), 20)
	name := ctx.FormValue("name")
	orderBy := ctx.FormValue("orderBy")

	roles := MGetAllRoles(name, orderBy, offset, limit)

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(apiResource(true, roles, "操作成功"))
}

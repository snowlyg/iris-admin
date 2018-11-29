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

	cp := t.ParseInt(ctx.FormValue("cp"), 1)
	mp := t.ParseInt(ctx.FormValue("mp"), 20)
	kw := ctx.FormValue("kw")
	roles := MGetAllRoles(kw, cp, mp)

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(apiResource(true, roles, "操作成功"))
}

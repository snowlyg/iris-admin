package controllers

import (
	"IrisYouQiKangApi/models"
	"IrisYouQiKangApi/system"
	"github.com/kataras/iris"
)

/**
* @api {get} /perms 获取所有的权限
* @apiName 获取所有的权限
* @apiGroup Perms
* @apiVersion 1.0.0
* @apiDescription 获取所有的权限
* @apiSampleRequest /perms
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func GetAllPerms(ctx iris.Context) {
	cp := system.Tools.ParseInt(ctx.FormValue("cp"), 1)
	mp := system.Tools.ParseInt(ctx.FormValue("mp"), 20)
	kw := ctx.FormValue("kw")
	perms := models.GetAllPerms(kw, cp, mp)

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(apiResource(true, perms, "操作成功"))
}

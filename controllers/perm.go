package controllers

import (
	"IrisYouQiKangApi/models"
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
	cp := Tools.ParseInt(ctx.FormValue("cp"), 1)
	mp := Tools.ParseInt(ctx.FormValue("mp"), 20)
	kw := ctx.FormValue("kw")

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(models.GetAllPerms(kw, cp, mp))
}

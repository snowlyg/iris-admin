package controllers

import (
	"IrisYouQiKangApi/models"
	"IrisYouQiKangApi/system"
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
func GetAllCompanies(ctx iris.Context) {
	cp := system.Tools.ParseInt(ctx.FormValue("cp"), 1)
	mp := system.Tools.ParseInt(ctx.FormValue("mp"), 20)
	kw := ctx.FormValue("kw")

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(models.GetAllCompanies(kw, cp, mp))
}

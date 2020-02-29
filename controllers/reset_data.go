package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/models"
	"github.com/snowlyg/IrisAdminApi/routepath"
)

/*
	重置系统数据
	管理端 管理员数据
	商户端 账号，角色，权限
*/
func ResetData(ctx iris.Context) {
	models.DelAllData()
	routes := routepath.GetRoutes(ctx.Application().GetRoutesReadOnly())
	models.CreateSystemData(routes) // 初始化系统数据 账号，角色，权限
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, routes, "重置数据成功"))
}

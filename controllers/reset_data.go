package controllers

import (
	"IrisAdminApi/config"
	"IrisAdminApi/models"
	"IrisAdminApi/routepath"
	"github.com/kataras/iris/v12"
)

func ResetData(ctx iris.Context) {
	models.DelAllData()
	routes := routepath.GetRoutes(ctx.Application().GetRoutesReadOnly())
	models.CreateSystemData(config.GetTfConf(), routes) // 初始化系统数据 管理员账号，角色，权限
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, routes, "重置数据成功"))
}

package controllers

import (
	"strings"

	"IrisAdminApi/config"
	"IrisAdminApi/models"
	"IrisAdminApi/validates"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func ResetData(ctx iris.Context) {
	models.DelAllData()
	routesReadOnly := getRoutes(ctx.Application().GetRoutesReadOnly())
	models.CreateSystemData(config.GetTfConf(), routesReadOnly) // 初始化系统数据 管理员账号，角色，权限
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, routesReadOnly, "重置数据成功"))
}

// 获取路由信息
func getRoutes(rs []context.RouteReadOnly) []*validates.PermissionRequest {
	//rs := api.APIBuilder.GetRoutes()
	var rrs []*validates.PermissionRequest
	for _, s := range rs {
		if !isPermRoute(s) {
			path := strings.Replace(s.Path(), ":id", "*", 1)
			rr := &validates.PermissionRequest{Name: path, DisplayName: s.Name(), Description: s.Name(), Act: s.Method()}
			rrs = append(rrs, rr)
		}
	}
	return rrs
}

// 过滤非必要权限
func isPermRoute(s context.RouteReadOnly) bool {
	exceptRouteName := []string{"OPTIONS", "GET", "POST", "HEAD", "PUT", "PATCH"}
	for _, er := range exceptRouteName {
		if strings.Contains(s.Name(), er) {
			return true
		}
	}
	return false
}

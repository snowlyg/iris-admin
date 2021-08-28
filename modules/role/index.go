package role

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/web"
)

// Party 调试模块
func Party() web.WebModule {
	handler := func(index iris.Party) {
		// index.Use(middleware.JwtHandler().Serve, middleware.New().ServeHTTP, middleware.OperationRecord())
		index.Get("/", GetAllRoles).Name = "角色列表"
		index.Get("/{id:uint}", GetRole).Name = "角色详情"
		index.Post("/", CreateRole).Name = "创建角色"
		index.Post("/{id:uint}", UpdateRole).Name = "编辑角色"
		index.Delete("/{id:uint}", DeleteRole).Name = "删除角色"
	}
	return web.NewModule("/users", handler)
}

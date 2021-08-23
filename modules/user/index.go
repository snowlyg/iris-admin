package user

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/web"
)

// Party 调试模块
func Party() web.WebModule {
	handler := func(index iris.Party) {
		// index.Use(middleware.JwtHandler().Serve, middleware.New().ServeHTTP, middleware.OperationRecord())
		// index.Get("/", GetUsers).Name = "用户列表"
		// index.Get("/{id:uint}", GetUser).Name = "用户详情"
		// index.Post("/", CreateUser).Name = "创建用户"
		// index.Post("/{id:uint}", UpdateUser).Name = "编辑用户"
		// index.Delete("/{id:uint}", DeleteUser).Name = "删除用户"
	}
	return web.NewModule("/users", handler)
}

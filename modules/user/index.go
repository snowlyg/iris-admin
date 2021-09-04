package user

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/module"
)

// Party 调试模块
func Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(),middleware.Casbin())
		index.Get("/", GetAllUsers).Name = "用户列表"
		index.Get("/{id:uint}", GetUser).Name = "用户详情"
		index.Post("/", CreateUser).Name = "创建用户"
		index.Post("/{id:uint}", UpdateUser).Name = "编辑用户"
		index.Delete("/{id:uint}", DeleteUser).Name = "删除用户"
		index.Get("/logout", Logout).Name = "退出"
		index.Get("/clear", Clear).Name = "清空 token"
		// index.Get("/expire", controllers.Expire).Name = "刷新 token"
		// index.Get("/profile", controllers.Profile).Name = "个人信息"
		// index.Post("/change_avatar", controllers.ChangeAvatar).Name = "修改头像"
	}
	return module.NewModule("/users", handler)
}

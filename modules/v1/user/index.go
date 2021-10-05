package user

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/module"
	"github.com/snowlyg/iris-admin/server/operation"
)

// Party 用户
func Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), operation.OperationRecord(), middleware.Casbin())
		index.Get("/", GetAll).Name = "用户列表"
		index.Get("/{id:uint}", GetUser).Name = "用户详情"
		index.Post("/", CreateUser).Name = "创建用户"
		index.Post("/{id:uint}", UpdateUser).Name = "编辑用户"
		index.Delete("/{id:uint}", DeleteUser).Name = "删除用户"
		index.Get("/logout", Logout).Name = "退出"
		index.Get("/clear", Clear).Name = "清空 token"
		index.Get("/profile", Profile).Name = "个人信息"
		index.Post("/change_avatar", ChangeAvatar).Name = "修改头像"
		// index.Get("/expire", controllers.Expire).Name = "刷新 token"
	}
	return module.NewModule("/users", handler)
}

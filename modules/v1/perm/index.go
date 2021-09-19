package perm

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/module"
)

// Party 权限
func Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())
		index.Get("/", GetAllPerms).Name = "权限列表"
		index.Get("/{id:uint}", GetPerm).Name = "权限详情"
		index.Post("/", CreatePerm).Name = "创建权限"
		index.Post("/{id:uint}", UpdatePerm).Name = "编辑权限"
		index.Delete("/{id:uint}", DeletePerm).Name = "删除权限"
	}
	return module.NewModule("/perms", handler)
}

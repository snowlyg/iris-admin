package perm

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/web"
)

// Party 调试模块
func Party() web.WebModule {
	handler := func(index iris.Party) {
		index.Get("/", GetAllPerms).Name = "权限列表"
		index.Get("/{id:uint}", GetPerm).Name = "权限详情"
		index.Post("/", CreatePerm).Name = "创建权限"
		index.Post("/{id:uint}", UpdatePerm).Name = "编辑权限"
		index.Delete("/{id:uint}", DeletePerm).Name = "删除权限"
	}
	return web.NewModule("/perms", handler)
}

package role

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/iris-admin/server/operation"
)

// Party 角色模块
func Party() func(index iris.Party) {
	return func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.MultiHandler(), operation.OperationRecord(), casbin.Casbin())
		index.Get("/", GetAll).Name = "角色列表"
		index.Get("/{id:uint}", First).Name = "角色详情"
		index.Post("/", CreateRole).Name = "创建角色"
		index.Post("/{id:uint}", UpdateRole).Name = "编辑角色"
		index.Delete("/{id:uint}", DeleteRole).Name = "删除角色"
	}
}

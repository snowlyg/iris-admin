package auth

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/iris-admin/server/operation"
)

// Party 认证模块
func Party() func(public iris.Party) {
	return func(public iris.Party) {
		public.Post("/login", Login)
		public.Use(middleware.MultiHandler(), casbin.Casbin(), operation.OperationRecord())
	}
}

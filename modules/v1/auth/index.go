package auth

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/module"
)

// Party 认证模块
func Party() module.WebModule {
	handler := func(public iris.Party) {
		public.Use(middleware.InitCheck())
		public.Post("/login", Login)
		public.Use(middleware.JwtHandler(), middleware.Casbin(), middleware.OperationRecord())
	}
	return module.NewModule("/auth", handler)
}

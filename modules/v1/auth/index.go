package auth

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/module"
)

// Party
func Party() module.WebModule {
	handler := func(public iris.Party) {
		public.Post("/admin/login", Login)
		public.Use(middleware.JwtHandler(), middleware.Casbin(), middleware.OperationRecord())
		// public.Post("/upload_file", iris.LimitRequestBodySize(libs.Config.MaxSize+1<<20), UploadFile).Name = "上传文件"
	}
	return module.NewModule("/auth", handler)
}

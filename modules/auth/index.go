package auth

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/module"
)

// Party
func Party() module.WebModule {
	handler := func(public iris.Party) {
		// 	v1.Post("/admin/login", controllers.Login)
		// index.Use(middleware.JwtHandler().Serve, middleware.New().ServeHTTP, middleware.OperationRecord())
		// 		index.Use(middleware.JwtHandler().Serve, middleware.New().ServeHTTP, middleware.OperationRecord()) //登录验证
		// index.Post("/upload_file", iris.LimitRequestBodySize(libs.Config.MaxSize+1<<20), controllers.UploadFile).Name = "上传文件"
	}
	return module.NewModule("/users", handler)
}

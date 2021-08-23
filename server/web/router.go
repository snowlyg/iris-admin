package web

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/middleware"
)

func (ws *WebServer) InitRouter() {
	ws.app.UseRouter(middleware.CrsAuth())
	app := ws.app.Party("/").AllowMethods(iris.MethodOptions)
	{

		for _, module := range ws.modules {
			app.PartyFunc(module.relativePath, module.handler)
		}

		// app.HandleDir("/uploads", iris.Dir(filepath.Join(libs.CWD(), "uploads")))
		// v1 := app.Party("api/v1")
		// {
		// 	// 是否开启接口请求频率限制
		// 	if !libs.Config.Limit.Disable {
		// 		limitV1 := rate.Limit(libs.Config.Limit.Limit, libs.Config.Limit.Burst, rate.PurgeEvery(time.Minute, 5*time.Minute))
		// 		v1.Use(limitV1)
		// 	}
		// 	v1.Post("/admin/login", controllers.Login)
		// 	v1.PartyFunc("/admin", func(admin iris.Party) { //casbin for gorm                                                   // <- IMPORTANT, register the middleware.
		// 		admin.Use(middleware.JwtHandler().Serve, middleware.New().ServeHTTP, middleware.OperationRecord()) //登录验证
		// 		admin.Get("/logout", controllers.Logout).Name = "退出"
		// 		admin.Get("/expire", controllers.Expire).Name = "刷新 token"
		// 		admin.Get("/clear", controllers.Clear).Name = "清空 token"
		// 		admin.Get("/profile", controllers.Profile).Name = "个人信息"
		// 		admin.Post("/change_avatar", controllers.ChangeAvatar).Name = "修改头像"
		// 		admin.Post("/upload_file", iris.LimitRequestBodySize(libs.Config.MaxSize+1<<20), controllers.UploadFile).Name = "上传文件"
		// 	})
		// }
	}
}

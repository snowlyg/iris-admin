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
		// }
	}
}

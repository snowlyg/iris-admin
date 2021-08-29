package web

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/middleware"
)

func (ws *WebServer) InitRouter() {
	ws.app.UseRouter(middleware.CrsAuth())
	app := ws.app.Party("/").AllowMethods(iris.MethodOptions)
	{
		fmt.Printf("config:%+v", g.CONFIG)
		if g.CONFIG.System.Level == "debug" {
			debug := DebugParty()
			app.PartyFunc(debug.relativePath, debug.handler)
		}

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

// Party 调试模块
func DebugParty() WebModule {
	handler := func(index iris.Party) {
		index.Get("/", func(ctx iris.Context) {
			ctx.HTML("<h1>请点击<a href='/debug/pprof'>这里</a>打开调试页面")
		})
		index.Any("/pprof", pprof.New())
		index.Any("/pprof/{action:path}", pprof.New())
	}
	return NewModule("/debug", handler)
}

package web

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/module"
)

func (ws *WebServer) InitRouter() {
	ws.app.UseRouter(middleware.CrsAuth())
	app := ws.app.Party("/").AllowMethods(iris.MethodOptions)
	{
		app.Use(middleware.InitCheck())
		if g.CONFIG.System.Level == "debug" {
			debug := DebugParty()
			app.PartyFunc(debug.RelativePath, debug.Handler)
		}
		ws.initModule()
		// app.HandleDir("/uploads", iris.Dir(filepath.Join(libs.CWD(), "uploads")))
	}
}

func (ws *WebServer) initModule() {
	fmt.Printf("modules : %d \n", len(ws.modules))
	if len(ws.modules) > 0 {
		for _, module := range ws.modules {
			sub := ws.app.PartyFunc(module.RelativePath, module.Handler)
			if len(module.Modules) > 0 {
				sub.PartyFunc(module.RelativePath, module.Handler)
			}
		}
	}
}

// Party 调试模块
func DebugParty() module.WebModule {
	handler := func(index iris.Party) {
		index.Get("/", func(ctx iris.Context) {
			ctx.HTML("<h1>请点击<a href='/debug/pprof'>这里</a>打开调试页面")
		})
		index.Any("/pprof", pprof.New())
		index.Any("/pprof/{action:path}", pprof.New())
	}
	return module.NewModule("/debug", handler)
}

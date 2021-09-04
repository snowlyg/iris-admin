package web

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/snowlyg/helper/arr"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/module"
	_ "github.com/snowlyg/iris-admin/server/viper"
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
		ws.GetSources()
	}
}

func (ws *WebServer) GetSources() {
	for _, r := range ws.app.GetRoutes() {
		// 去除非接口路径
		handerNames := context.HandlersNames(r.Handlers)
		if !arr.InArrayS([]string{"GET", "POST", "PUT", "DELETE"}, r.Method) || r.IsStatic() || !arr.InArrayS(strings.Split(handerNames, ","), "github.com/snowlyg/multi.(*Verifier).Verify") {
			continue
		}
		ws.wg.Add(1)
		go func(r *router.Route) {
			g.PermRoutes = append(g.PermRoutes, map[string]string{
				"path": r.Path,
				"name": r.Name,
				"act":  r.Method,
			})
			ws.wg.Done()
		}(r)
		ws.wg.Wait()
	}
}

func (ws *WebServer) initModule() {
	if len(ws.modules) > 0 {
		for _, mod := range ws.modules {
			ws.wg.Add(1)
			go func(mod module.WebModule) {
				sub := ws.app.PartyFunc(mod.RelativePath, mod.Handler)
				if len(mod.Modules) > 0 {
					for _, subModule := range mod.Modules {
						sub.PartyFunc(subModule.RelativePath, subModule.Handler)
					}
				}
				ws.wg.Done()
			}(mod)
		}
		ws.wg.Wait()
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

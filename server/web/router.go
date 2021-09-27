package web

import (
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/snowlyg/helper/arr"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/config"
	"github.com/snowlyg/iris-admin/server/module"
)

// InitRouter 初始化模块路由
func (ws *WebServer) InitRouter() error {
	ws.app.UseRouter(middleware.CrsAuth())

	app := ws.app.Party("/").AllowMethods(iris.MethodOptions)
	{
		app.Use(middleware.InitCheck())
		if config.CONFIG.System.Level == "debug" {
			debug := DebugParty()
			app.PartyFunc(debug.RelativePath, debug.Handler)
		}
		ws.initModule()
		ws.AddUploadStatic()
		ws.AddWebStatic("/")
		err := ws.app.Build()
		if err != nil {
			return fmt.Errorf("build router %w", err)
		}
		g.PermRoutes = ws.GetSources()
		return nil
	}
}

// GetSources 获取web服务需要认证的权限
func (ws *WebServer) GetSources() []map[string]string {
	routeLen := len(ws.app.GetRoutes())
	ch := make(chan map[string]string, routeLen)
	for _, r := range ws.app.GetRoutes() {
		r := r
		// 去除非接口路径
		handerNames := context.HandlersNames(r.Handlers)
		if !arr.InArrayS([]string{"GET", "POST", "PUT", "DELETE"}, r.Method) || !arr.InArrayS(strings.Split(handerNames, ","), "github.com/snowlyg/multi.(*Verifier).Verify") {
			routeLen--
			continue
		}
		go func(r *router.Route) {
			route := map[string]string{
				"path": r.Path,
				"name": r.Name,
				"act":  r.Method,
			}
			ch <- route
		}(r)
	}

	routes := make([]map[string]string, routeLen)
	for i := 0; i < routeLen; i++ {
		routes[i] = <-ch
	}
	return routes
}

// initModule 初始化web服务模块，包括子模块
func (ws *WebServer) initModule() {
	if len(ws.modules) > 0 {
		for _, mod := range ws.modules {
			mod := mod
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

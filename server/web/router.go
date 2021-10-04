package web

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/snowlyg/helper/arr"
	"github.com/snowlyg/iris-admin/middleware"
	"github.com/snowlyg/iris-admin/server/module"
)

var (
	PermRoutes   chan map[string]string
	NoPermRoutes chan map[string]string
)

// InitRouter 初始化模块路由
func (ws *WebServer) InitRouter() error {
	app := ws.app.Party("/").AllowMethods(iris.MethodOptions)
	{
		app.UseRouter(middleware.CrsAuth())
		app.Use(middleware.InitCheck())
		if CONFIG.System.Level == "debug" {
			debug := DebugParty()
			app.PartyFunc(debug.RelativePath, debug.Handler)
		}
		// app.PartyFunc("/init", InitParty().Handler)
	}
	ws.initModule()
	ws.AddUploadStatic()
	ws.AddWebStatic("/")
	ws.GetSources()
	return nil
}

// GetSources 获取系统路由
// - perm 权鉴路由
// - noPerm 公共路由
func (ws *WebServer) GetSources() {
	routeLen := len(ws.app.GetRoutes())
	PermRoutes = make(chan map[string]string, routeLen)
	NoPermRoutes = make(chan map[string]string, routeLen)
	for _, r := range ws.app.GetRoutes() {
		r := r
		go func(r *router.Route) {
			route := map[string]string{
				"path": r.Path,
				"name": r.Name,
				"act":  r.Method,
			}
			handerNames := context.HandlersNames(r.Handlers)
			if !arr.InArrayS([]string{"GET", "POST", "PUT", "DELETE"}, r.Method) || !arr.InArrayS(strings.Split(handerNames, ","), "github.com/snowlyg/multi.(*Verifier).Verify") {
				NoPermRoutes <- route
			} else {
				PermRoutes <- route
			}

		}(r)
	}
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

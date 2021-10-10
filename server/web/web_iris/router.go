package web_iris

import (
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/kataras/iris/v12/middleware/rate"
	"github.com/snowlyg/helper/arr"
	"github.com/snowlyg/iris-admin/middleware"
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
		if !CONFIG.Limit.Disable {
			limitV1 := rate.Limit(CONFIG.Limit.Limit, CONFIG.Limit.Burst, rate.PurgeEvery(time.Minute, 5*time.Minute))
			app.Use(limitV1)
		}
		if CONFIG.System.Level == "debug" {
			debug := func(index iris.Party) {
				index.Get("/", func(ctx iris.Context) {
					ctx.HTML("<h1>请点击<a href='/debug/pprof'>这里</a>打开调试页面")
				})
				index.Any("/pprof", pprof.New())
				index.Any("/pprof/{action:path}", pprof.New())
			}
			app.PartyFunc("/debug", debug)
		}

		for _, party := range ws.parties {
			app.PartyFunc(party.Perfix, party.PartyFunc)
		}
		// app.PartyFunc("/init", InitParty().Handler)
	}
	if ws.staticPrefix != "" {
		ws.AddUploadStatic()
	}
	if ws.webPrefix != "" {
		ws.AddWebStatic(ws.webPrefix)
	}
	ws.GetSources()
	return nil
}

// GetSources 获取系统路由
// - PermRoutes 权鉴路由
// - NoPermRoutes 公共路由
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

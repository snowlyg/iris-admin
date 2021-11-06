package web_gin

import (
	limit "github.com/aviddiviner/gin-limit"
	"github.com/snowlyg/iris-admin/server/web/web_gin/middleware"
)

// InitRouter 初始化模块路由
func (ws *WebServer) InitRouter() error {
	ws.app.Use(limit.MaxAllowed(50))
	app := ws.app.Group("/")
	{
		app.Use(middleware.Cors()) // 如需跨域可以打开
		// if CONFIG.System.Level == "debug" {
		// 	debug := func(index iris.Party) {
		// 		index.Get("/", func(ctx iris.Context) {
		// 			ctx.HTML("<h1>请点击<a href='/debug/pprof'>这里</a>打开调试页面")
		// 		})
		// 		index.Any("/pprof", pprof.New())
		// 		index.Any("/pprof/{action:path}", pprof.New())
		// 	}
		// 	app.PartyFunc("/debug", debug)
		// }

		// for _, party := range ws.parties {
		// 	app.PartyFunc(party.Perfix, party.PartyFunc)
		// }
	}
	if ws.staticPrefix != "" {
		ws.AddUploadStatic()
	}
	if ws.webPrefix != "" {
		ws.AddWebStatic()
	}
	return nil
}

// GetSources 获取系统路由
// - PermRoutes 权鉴路由
// - NoPermRoutes 公共路由
func (ws *WebServer) GetSources() ([]map[string]string, []map[string]string) {
	routeLen := len(ws.app.Routes())
	permRoutes := make([]map[string]string, 0, routeLen)
	noPermRoutes := make([]map[string]string, 0, routeLen)
	for _, r := range ws.app.Routes() {

		route := map[string]string{
			"path": r.Path,
			"name": r.Handler,
			"act":  r.Method,
		}
		// handerNames := context.HandlersNames(r.HandlerFunc)
		// if !arr.InArrayS([]string{"GET", "POST", "PUT", "DELETE"}, r.Method) || !arr.InArrayS(strings.Split(handerNames, ","), "github.com/snowlyg/multi.(*Verifier).Verify") {
		// 	noPermRoutes = append(noPermRoutes, route)
		// } else {
		permRoutes = append(permRoutes, route)
		// }
	}
	return permRoutes, noPermRoutes
}

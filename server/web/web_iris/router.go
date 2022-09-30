package web_iris

import (
	"net/http"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/kataras/iris/v12/middleware/rate"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/snowlyg/helper/arr"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_iris/middleware"
)

// InitRouter
func (ws *WebServer) InitRouter() error {
	app := ws.app.Party("/").AllowMethods(iris.MethodOptions)
	{
		app.Get("/v0/version", func(ctx iris.Context) {
			ctx.WriteString("IRIS-ADMIN is running!!!")
		})

		app.UseRouter(middleware.CrsAuth())
		app.UseRouter(recover.New())
		if !web.CONFIG.Limit.Disable {
			limitV1 := rate.Limit(web.CONFIG.Limit.Limit, web.CONFIG.Limit.Burst, rate.PurgeEvery(time.Minute, 5*time.Minute))
			app.Use(limitV1)
		}
		if web.CONFIG.System.Level == "debug" {
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
	}

	// http test must build
	if err := ws.app.Build(); err != nil {
		return err
	}

	return nil
}

// GetSources
// - PermRoutes
// - NoPermRoutes
func (ws *WebServer) GetSources() ([]map[string]string, []map[string]string) {
	methodExcepts := strings.Split(web.CONFIG.Except.Method, ";")
	uris := strings.Split(web.CONFIG.Except.Uri, ";")
	routeLen := len(ws.app.GetRoutes())
	permRoutes := make([]map[string]string, 0, routeLen)
	noPermRoutes := make([]map[string]string, 0, routeLen)

	for _, r := range ws.app.GetRoutes() {
		route := map[string]string{
			"path": r.Path,
			"name": r.Name,
			"act":  r.Method,
		}

		if !arr.InArray([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}, r.Method) {
			noPermRoutes = append(noPermRoutes, route)
			continue
		}

		if len(methodExcepts) > 0 && len(uris) > 0 && len(methodExcepts) == len(uris) {
			for i := 0; i < len(methodExcepts); i++ {
				if strings.EqualFold(r.Method, strings.ToLower(methodExcepts[i])) && strings.EqualFold(r.Path, strings.ToLower(uris[i])) {
					noPermRoutes = append(noPermRoutes, route)
					continue
				}
			}
		}

		permRoutes = append(permRoutes, route)
	}
	return permRoutes, noPermRoutes
}

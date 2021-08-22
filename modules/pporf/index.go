package pporf

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/snowlyg/iris-admin/server/web"
)

// Party 返回PPROF模块
func Party() web.WebModule {
	handler := func(index iris.Party) {
		index.Get("/", func(ctx iris.Context) {
			ctx.HTML("<h1>请点击<a href='/debug/pprof'>这里</a>进去调试页面")
		})
		index.Any("/debug/pprof", pprof.New())
		index.Any("/debug/pprof/{action:path}", pprof.New())
	}
	return web.NewModule("/debug", handler)
}

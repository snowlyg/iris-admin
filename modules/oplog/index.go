package oplog

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/web"
)

// Party 调试模块
func Party() web.WebModule {
	handler := func(index iris.Party) {
		// index.Use(middleware.JwtHandler().Serve, middleware.New().ServeHTTP)
	}
	return web.NewModule("/oplog", handler)
}

package middleware

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/config"
	"github.com/snowlyg/iris-admin/server/database"
)

// InitCheck 初始化检测中间件
func InitCheck() iris.Handler {
	return func(ctx *context.Context) {
		if database.Instance() == nil || (config.CONFIG.System.CacheType == "redis" && cache.Instance() == nil) {
			ctx.StopWithJSON(http.StatusOK, g.Response{Code: g.NeedInitErr.Code, Data: nil, Msg: g.NeedInitErr.Msg})
		} else {
			ctx.Next()
		}
		// 处理请求
	}
}

package middleware

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
)

// InitCheck 初始化检测中间件
func InitCheck() iris.Handler {
	return func(ctx *context.Context) {
		if database.Instance() == nil || cache.Instance() == nil {
			ctx.StopWithJSON(http.StatusOK, orm.Response{Code: orm.NeedInitErr.Code, Data: nil, Msg: orm.NeedInitErr.Msg})
		} else {
			ctx.Next()
		}
		// 处理请求
	}
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/iris-admin/server/web/web_gin/response"
	"github.com/snowlyg/iris-admin/server/zap_server"
	multi "github.com/snowlyg/multi/gin"
)

// 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		obj := ctx.Request.URL.RequestURI()
		// 获取请求方法
		act := ctx.Request.Method
		// 获取用户的角色
		sub := multi.GetAuthorityId(ctx)

		if sub == "" {
			zap_server.ZAPLOG.Info("user authorityId is empty")
			response.UnauthorizedFailWithMessage("auth token 已经过期", ctx)
			ctx.Abort()
			return
		}

		success, err := casbin.Instance().Enforce(sub, obj, act)
		if err != nil {
			response.ForbiddenFailWithMessage("权限服务验证失败：verfiy failed", ctx)
			ctx.Abort()
			return
		}
		if !success {
			response.ForbiddenFailWithMessage("无此操作权限", ctx)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

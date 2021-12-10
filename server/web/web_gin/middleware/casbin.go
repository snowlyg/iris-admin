package middleware

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/iris-admin/server/web/web_gin/response"
	"github.com/snowlyg/iris-admin/server/zap_server"
	multi "github.com/snowlyg/multi/gin"
)

// 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		obj := filepath.ToSlash(filepath.Clean(ctx.Request.URL.Path))
		// 获取请求方法
		act := ctx.Request.Method
		// 获取用户的角色
		subs := multi.GetAuthorityId(ctx)
		if len(subs) == 0 {
			zap_server.ZAPLOG.Info("用户角色ID为空")
			response.UnauthorizedFailWithMessage("TOKEN已经过期", ctx)
			ctx.Abort()
			return
		}

		for _, sub := range subs {
			success, err := casbin.Instance().Enforce(sub, obj, act)
			if err != nil {
				response.ForbiddenFailWithMessage("权限服务验证失败", ctx)
				continue
			}
			if success {
				ctx.Next()
				return
			}
		}

		response.ForbiddenFailWithMessage("无此操作权限", ctx)
		ctx.Abort()
	}
}

package web_gin

import (
	"net/http"
	"strings"

	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/helper/arr"
	"github.com/snowlyg/iris-admin/server/web/web_gin/middleware"
)

type CustomAspect struct {
	CustomValue int
}

func (a *CustomAspect) GetStats() interface{} {
	return a.CustomValue
}

func (a *CustomAspect) Name() string {
	return "Custom"
}

func (a *CustomAspect) InRoot() bool {
	return false
}

func (ws *WebServer) GetRouterGroup(relativePath string) *gin.RouterGroup {
	return ws.app.Group(relativePath)
}

// InitRouter 初始化模块路由
func (ws *WebServer) InitRouter() error {
	ws.app.Use(limit.MaxAllowed(50))
	router := ws.app.Group("/")
	{
		router.Use(middleware.Cors()) // 如需跨域可以打开
		// last middleware
		router.Use(gin.Recovery())

		// 排除路由竞争
		if ws.webPrefix != "/" {
			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "GO_MERCHANT is running!!!")
			})
		}
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
			"path":   r.Path,
			"name":   "",
			"method": r.Method,
		}

		if !arr.InArrayS([]string{"GET", "POST", "PUT", "DELETE"}, r.Method) || !strings.Contains(r.Handler, "github.com/snowlyg/multi.(*Verifier).Verify") {
			noPermRoutes = append(noPermRoutes, route)
		} else {
			permRoutes = append(permRoutes, route)
		}
	}
	return permRoutes, noPermRoutes
}

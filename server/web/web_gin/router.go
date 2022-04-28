package web_gin

import (
	"net/http"
	"path/filepath"
	"strings"

	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/helper/arr"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_gin/middleware"
)

func (ws *WebServer) GetRouterGroup(relativePath string) *gin.RouterGroup {
	return ws.app.Group(relativePath)
}

// InitRouter 初始化模块路由
func (ws *WebServer) InitRouter() error {
	ws.app.Use(limit.MaxAllowed(50))
	ws.app.Use(gin.Recovery())
	if web.CONFIG.System.Level == "debug" {
		pprof.Register(ws.app)
	}
	router := ws.app.Group("/")
	{
		router.Use(middleware.Cors()) // 如需跨域可以打开
		// last middleware
		router.Use(gin.Recovery())

		// 排除路由竞争
		router.GET("/v0/version", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "IRIS-ADMIN is running!!!")
		})
	}
	return nil
}

// GetSources 获取系统路由
// - PermRoutes 权鉴路由
// - NoPermRoutes 公共路由
func (ws *WebServer) GetSources() ([]map[string]string, []map[string]string) {

	methodExcepts := strings.Split(web.CONFIG.Except.Method, ";")
	uriExcepts := strings.Split(web.CONFIG.Except.Uri, ";")

	routeLen := len(ws.app.Routes())
	permRoutes := make([]map[string]string, 0, routeLen)
	noPermRoutes := make([]map[string]string, 0, routeLen)

	for _, r := range ws.app.Routes() {
		// 处理路径
		bases := strings.Split(filepath.Base(r.Handler), ".")
		if len(bases) != 2 {
			continue
		}
		path := filepath.ToSlash(filepath.Clean(r.Path))

		route := map[string]string{
			"path":   path,
			"desc":   bases[1],
			"group":  bases[0],
			"method": r.Method,
		}
		// 过滤非必要请求
		if !arr.InArrayS([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}, r.Method) {
			noPermRoutes = append(noPermRoutes, route)
			continue
		}

		// 过滤不需要的请求
		if len(methodExcepts) > 0 && len(uriExcepts) > 0 && len(methodExcepts) == len(uriExcepts) {
			for i := 0; i < len(methodExcepts); i++ {
				if strings.EqualFold(r.Method, strings.ToLower(methodExcepts[i])) && strings.EqualFold(path, strings.ToLower(uriExcepts[i])) {
					noPermRoutes = append(noPermRoutes, route)
					continue
				}
			}
		}

		permRoutes = append(permRoutes, route)
	}
	return permRoutes, noPermRoutes
}

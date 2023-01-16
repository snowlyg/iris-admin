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

// InitRouter
func (ws *WebServer) InitRouter() error {
	ws.app.Use(limit.MaxAllowed(50))
	ws.app.Use(gin.Recovery())
	if web.CONFIG.System.Level == "debug" {
		pprof.Register(ws.app)
	}
	router := ws.app.Group("/")
	{
		router.Use(middleware.Cors())
		// last middleware
		router.Use(gin.Recovery())

		router.GET("/v0/version", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "IRIS-ADMIN is running!!!")
		})
	}
	return nil
}

// GetSources
// - PermRoutes
// - NoPermRoutes
func (ws *WebServer) GetSources() ([]map[string]string, []map[string]string) {

	methodExcepts := strings.Split(web.CONFIG.Except.Method, ";")
	uriExcepts := strings.Split(web.CONFIG.Except.Uri, ";")

	routeLen := len(ws.app.Routes())
	permRoutes := make([]map[string]string, 0, routeLen)
	noPermRoutes := make([]map[string]string, 0, routeLen)

	for _, r := range ws.app.Routes() {
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
		httpStatusType := arr.NewCheckArrayType(4)
		httpStatusType.AddMutil(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
		if !httpStatusType.Check(r.Method) {
			noPermRoutes = append(noPermRoutes, route)
			continue
		}

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

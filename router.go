package admin

import (
	"net/http"
	"path/filepath"
	"strings"

	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/helper/arr"
)

func (ws *WebServe) GetRouterGroup(relativePath string) *gin.RouterGroup {
	return ws.app.Group(relativePath)
}

// InitRouter
func (ws *WebServe) InitRouter() error {
	ws.app.Use(limit.MaxAllowed(50))
	if ws.conf.System.Level == "debug" {
		pprof.Register(ws.app)
	}
	router := ws.app.Group("/")
	{
		// router.Use(Cors())

		router.GET("/v0/version", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "IRIS-ADMIN is running!!!")
		})
	}
	return nil
}

// GetSources
// - PermRoutes
// - NoPermRoutes
func (ws *WebServe) GetSources() ([]map[string]string, []map[string]string) {
	methodExcepts := strings.Split(ws.conf.Except.Method, ";")
	uriExcepts := strings.Split(ws.conf.Except.Uri, ";")
	methodMenus := strings.Split(ws.conf.Menu.Method, ";")
	uriMenus := strings.Split(ws.conf.Menu.Uri, ";")

	routeLen := len(ws.app.Routes())
	permRoutes := make([]map[string]string, 0, routeLen)
	otherMethodTypes := make([]map[string]string, 0, routeLen)

	for _, r := range ws.app.Routes() {
		bases := strings.Split(filepath.Base(r.Handler), ".")
		if len(bases) != 2 {
			continue
		}
		path := filepath.ToSlash(filepath.Clean(r.Path))
		route := map[string]string{
			"path":    path,
			"desc":    bases[1],
			"group":   bases[0],
			"method":  r.Method,
			"is_menu": "0",
		}
		if len(methodMenus) > 0 && len(uriMenus) > 0 && len(methodMenus) == len(uriMenus) {
			for i := 0; i < len(methodMenus); i++ {
				if strings.EqualFold(r.Method, strings.ToLower(methodMenus[i])) && strings.EqualFold(path, strings.ToLower(uriMenus[i])) {
					route["is_menu"] = "1"
				}
			}
		}

		httpStatusType := arr.NewCheckArrayType(4)
		httpStatusType.AddMutil(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
		if !httpStatusType.Check(r.Method) {
			otherMethodTypes = append(otherMethodTypes, route)
			continue
		}

		if len(methodExcepts) > 0 && len(uriExcepts) > 0 && len(methodExcepts) == len(uriExcepts) {
			for i := 0; i < len(methodExcepts); i++ {
				if strings.EqualFold(r.Method, strings.ToLower(methodExcepts[i])) && strings.EqualFold(path, strings.ToLower(uriExcepts[i])) {
					otherMethodTypes = append(otherMethodTypes, route)
					continue
				}
			}
		}
		permRoutes = append(permRoutes, route)
	}
	return permRoutes, otherMethodTypes
}

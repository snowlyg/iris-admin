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

func (ws *WebServe) Group(relativePath string) *gin.RouterGroup {
	return ws.engine.Group(relativePath)
}

func (ws *WebServe) InitRouter() error {
	ws.engine.Use(limit.MaxAllowed(50))
	if ws.conf.System.Level == "debug" {
		pprof.Register(ws.engine)
	}
	router := ws.engine.Group("/")
	{
		router.GET("/health", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "IRIS-ADMIN IS RUNNING!!!")
		})
	}
	return nil
}

func (ws *WebServe) GetSources() ([]map[string]string, []map[string]string) {
	methodExcepts := strings.Split(ws.conf.Except.Method, ";")
	uriExcepts := strings.Split(ws.conf.Except.Uri, ";")

	routeLen := len(ws.engine.Routes())
	permRoutes := make([]map[string]string, 0, routeLen)
	otherMethodTypes := make([]map[string]string, 0, routeLen)

	for _, r := range ws.engine.Routes() {
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

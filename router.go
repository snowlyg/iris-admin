package admin

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/helper/arr"
	"gorm.io/gorm"
)

type Router struct {
	gorm.Model
	Path     string    `json:"path"`
	Title    string    `json:"title"`
	Group    string    `json:"group"`
	Method   string    `json:"method"`
	Children []*Router `json:"children" gorm:"-"`
}

func (m *Router) TableName() string {
	return "routers"
}

type Menu struct {
	gorm.Model
	Path       string `json:"path"`
	Component  string `json:"component"`
	Redirect   string `json:"redirect"`
	Hidden     bool   `json:"hidden"`
	AlwaysShow bool   `json:"alwaysShow"`
	Meta
	Children []*Menu `json:"children" gorm:"-"`
}

func (m *Menu) TableName() string {
	return "menus"
}

type Meta struct {
	Roles   []string `json:"roles" gorm:"-"`
	Title   string   `json:"title"`
	Icon    string   `json:"icon"`
	NoCache bool     `json:"noCache"`
}

func (ws *WebServe) InitRouter() error {
	ws.engine.Use(limit.MaxAllowed(50))
	log.Printf("use gin-limit middleware\n")
	if ws.conf.System.Level == "debug" {
		pprof.Register(ws.engine)
	}
	ws.engine.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "IRIS-ADMIN IS RUNNING!!!")
	})
	return nil
}

func (ws *WebServe) Routers() {
	methodExcepts := strings.Split(ws.conf.Except.Method, ";")
	uriExcepts := strings.Split(ws.conf.Except.Uri, ";")

	// routeLen := len(ws.engine.Routes())
	// permRoutes := make([]*Router, 0, routeLen)
	// otherMethodTypes := make([]*Router, 0, routeLen)

	for _, r := range ws.engine.Routes() {
		bases := strings.Split(filepath.Base(r.Handler), ".")
		if len(bases) != 2 {
			continue
		}
		path := filepath.ToSlash(filepath.Clean(r.Path))
		route := &Router{
			Path:   path,
			Title:  bases[1],
			Group:  bases[0],
			Method: r.Method,
		}

		httpStatusType := arr.NewCheckArrayType(4)
		httpStatusType.AddMutil(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
		if !httpStatusType.Check(r.Method) {
			ws.otherRoutes = append(ws.otherRoutes, route)
			continue
		}

		if len(methodExcepts) > 0 && len(uriExcepts) > 0 && len(methodExcepts) == len(uriExcepts) {
			for i := 0; i < len(methodExcepts); i++ {
				if strings.EqualFold(r.Method, strings.ToLower(methodExcepts[i])) && strings.EqualFold(path, strings.ToLower(uriExcepts[i])) {
					ws.otherRoutes = append(ws.otherRoutes, route)
					continue
				}
			}
		}
		ws.permRoutes = append(ws.permRoutes, route)
	}

	if ws.db != nil {

	}
}

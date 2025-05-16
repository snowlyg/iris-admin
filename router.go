package admin

import (
	"log"
	"net/http"
	"strings"

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

func Group(group *gin.RouterGroup) {
	r := group.Group("/routes")
	{
		r.GET("/list", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"status":  200,
				"message": "OK",
			})
		})
	}
}

func (ws *WebServe) groupRouters() {
	methodExcepts := strings.Split(ws.conf.Except.Method, ";")
	uriExcepts := strings.Split(ws.conf.Except.Uri, ";")

	// routeLen := len(ws.engine.Routes())
	// permRoutes := make([]*Router, 0, routeLen)
	// otherMethodTypes := make([]*Router, 0, routeLen)

	for _, r := range ws.engine.Routes() {
		// log.Printf("handler:%s, method:%s, path:%s\n", r.Handler, r.Method, r.Path)
		if strings.Contains(r.Path, "/*filepath") || r.Handler == "github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1" {
			continue
		}
		path := r.Path
		route := &Router{
			Path:   path,
			Title:  path,
			Group:  "",
			Method: r.Method,
		}

		httpStatusType := arr.NewCheckArrayType(4)
		httpStatusType.AddMutil(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
		if !httpStatusType.Check(r.Method) {
			ws.otherRoutes = append(ws.otherRoutes, route)
			continue
		}

		if len(methodExcepts) > 0 && len(uriExcepts) > 0 && len(methodExcepts) == len(uriExcepts) {
			for i := range methodExcepts {
				if strings.EqualFold(r.Method, strings.ToLower(methodExcepts[i])) && strings.EqualFold(path, strings.ToLower(uriExcepts[i])) {
					ws.otherRoutes = append(ws.otherRoutes, route)
					continue
				}
			}
		}
		ws.permRoutes = append(ws.permRoutes, route)
	}

	// log.Printf("permRoutes:%d other:%d\n", len(ws.permRoutes), len(ws.otherRoutes))

	if ws.db == nil {
		return
	}

	if len(ws.permRoutes) == 0 {
		return
	}

	// seed routers
	olds := []*Router{}
	dels := []uint{}
	adds := []*Router{}
	if err := ws.db.Model(&Router{}).Find(&olds).Error; err != nil {
		log.Printf("iris-admin: old router find get err:%s\n", err.Error())
	}

	if len(olds) == 0 {
		if err := ws.db.Create(&ws.permRoutes).Error; err == nil {
			log.Printf("iris-admin: add %d router \n", len(ws.permRoutes))
		}
		return
	}

	oldCheck := arr.NewCheckArrayType(len(olds))
	for _, old := range olds {
		oldCheck.Add(old.Path)
		found := false
		for _, a := range ws.permRoutes {
			if old.Path == a.Path && old.Method == a.Method {
				found = true
				break
			}
		}
		if !found {
			dels = append(dels, old.ID)
		}
	}

	if len(dels) > 0 {
		if err := ws.db.Delete(&Router{}, dels).Error; err == nil {
			log.Printf("iris-admin: delete %d router\n", len(dels))
		}
	}

	for _, r := range ws.permRoutes {
		if !oldCheck.Check(r.Path) {
			adds = append(adds, r)
		}
	}

	if len(adds) > 0 {
		if err := ws.db.Create(&adds).Error; err == nil {
			log.Printf("iris-admin: add %d router,old:%d\n", len(adds), len(olds))
		}
	}
}

package admin

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
func (m *Router) List() []map[string]any {
	return []map[string]any{}
}

type routeKey struct {
	path   string
	method string
}

var permissionMethods = map[string]struct{}{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodPut:    {},
	http.MethodDelete: {},
}

func normalizeRoutePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" || strings.HasPrefix(path, "/") {
		return path
	}
	return "/" + path
}

func newRouteKey(path, method string) routeKey {
	return routeKey{
		path:   normalizeRoutePath(path),
		method: strings.ToUpper(strings.TrimSpace(method)),
	}
}

func routeExceptions(methods, paths string) map[routeKey]struct{} {
	methodItems := strings.Split(methods, ";")
	pathItems := strings.Split(paths, ";")
	exceptions := make(map[routeKey]struct{})
	if len(methodItems) != len(pathItems) {
		return exceptions
	}
	for i := range methodItems {
		key := newRouteKey(pathItems[i], methodItems[i])
		if key.path == "" || key.method == "" {
			continue
		}
		exceptions[key] = struct{}{}
	}
	return exceptions
}

func classifyRoutes(routes gin.RoutesInfo, exceptMethods, exceptPaths string) (permRoutes, otherRoutes []*Router) {
	exceptions := routeExceptions(exceptMethods, exceptPaths)
	for _, info := range routes {
		if strings.Contains(info.Path, "/*filepath") ||
			info.Handler == "github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1" {
			continue
		}
		key := newRouteKey(info.Path, info.Method)
		route := &Router{
			Path:   key.path,
			Title:  key.path,
			Method: key.method,
		}
		if _, ok := permissionMethods[key.method]; !ok {
			otherRoutes = append(otherRoutes, route)
			continue
		}
		if _, ok := exceptions[key]; ok {
			otherRoutes = append(otherRoutes, route)
			continue
		}
		permRoutes = append(permRoutes, route)
	}
	return permRoutes, otherRoutes
}

func diffRoutes(existing, desired []*Router) (deleteIDs []uint, additions []*Router) {
	existingKeys := make(map[routeKey]struct{}, len(existing))
	desiredKeys := make(map[routeKey]struct{}, len(desired))
	for _, route := range desired {
		desiredKeys[newRouteKey(route.Path, route.Method)] = struct{}{}
	}
	for _, route := range existing {
		key := newRouteKey(route.Path, route.Method)
		existingKeys[key] = struct{}{}
		if _, ok := desiredKeys[key]; !ok {
			deleteIDs = append(deleteIDs, route.ID)
		}
	}
	for _, route := range desired {
		if _, ok := existingKeys[newRouteKey(route.Path, route.Method)]; !ok {
			additions = append(additions, route)
		}
	}
	return deleteIDs, additions
}

func (ws *WebServe) groupRouters() {
	routes := ws.engine.Routes()
	for _, route := range routes {
		log.Printf("handler:%s, method:%s, path:%s\n", route.Handler, route.Method, route.Path)
	}
	ws.permRoutes, ws.otherRoutes = classifyRoutes(routes, ws.conf.Except.Method, ws.conf.Except.Uri)

	if ws.db == nil {
		return
	}

	olds := []*Router{}
	if err := ws.db.Model(&Router{}).Find(&olds).Error; err != nil {
		log.Printf("iris-admin: old router find get err:%s\n", err.Error())
		return
	}

	dels, adds := diffRoutes(olds, ws.permRoutes)
	if len(dels) > 0 {
		if err := ws.db.Delete(&Router{}, dels).Error; err != nil {
			log.Printf("iris-admin: delete routers failed:%s\n", err.Error())
		} else {
			log.Printf("iris-admin: delete %d routers\n", len(dels))
		}
	}

	if len(adds) > 0 {
		if err := ws.db.Create(&adds).Error; err != nil {
			log.Printf("iris-admin: add routers failed:%s\n", err.Error())
		} else {
			log.Printf("iris-admin: add %d routers,old:%d\n", len(adds), len(olds))
		}
	}
}

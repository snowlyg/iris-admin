package routepath

import (
	"strings"

	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/snowlyg/IrisAdminApi/validates"
)

type PathName struct {
	Name   string
	Path   string
	Method string
}

func getPathNames(i interface{}) []*PathName {
	var pns []*PathName
	if routeReadOnly, ok := i.([]context.RouteReadOnly); ok {
		for _, s := range routeReadOnly {
			pn := &PathName{
				Name:   s.Name(),
				Path:   s.Path(),
				Method: s.Method(),
			}
			pns = append(pns, pn)
		}
	} else if route, ok := i.([]*router.Route); ok {
		for _, s := range route {
			pn := &PathName{
				Name:   s.Name,
				Path:   s.Path,
				Method: s.Method,
			}
			pns = append(pns, pn)
		}
	}
	return pns
}

// 获取路由信息
func GetRoutes(i interface{}) []*validates.PermissionRequest {
	var rrs []*validates.PermissionRequest
	for _, s := range getPathNames(i) {
		if !isPermRoute(s.Name) {
			rr := &validates.PermissionRequest{Name: s.Path, DisplayName: s.Name, Description: s.Name, Act: s.Method}
			rrs = append(rrs, rr)
		}
	}
	return rrs
}

// 过滤非必要权限
func isPermRoute(name string) bool {
	exceptRouteName := []string{"OPTIONS", "GET", "POST", "HEAD", "PUT", "PATCH", "payload"}
	for _, er := range exceptRouteName {
		if strings.Contains(name, er) {
			return true
		}
	}
	return false
}

package app

import (
	stdContext "context"
	"fmt"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/models"
	"github.com/snowlyg/blog/routes"
	"github.com/snowlyg/easygorm"
)

// Server
type Server struct {
	App    *iris.Application
	Status bool // 服务状态
}

func NewServer() *Server {
	app := iris.New()

	app.Logger().SetLevel(libs.Config.LogLevel)                       //设置日志级别
	libs.InitRedisCluster(libs.GetRedisUris(), libs.Config.Redis.Pwd) //初始化redis
	models.Migrate()                                                  //初始化模型
	routes.App(app)                                                   //注册 app 路由

	// CTRL+C/CMD+C pressed or a unix kill command received
	iris.RegisterOnInterrupt(func() {
		db, err := easygorm.Egm.Db.DB()
		if err != nil {
			panic(err)
		}
		defer db.Close()
	})

	return &Server{
		App:    app,
		Status: false,
	}
}

// Start
func (s *Server) Start() error {
	if libs.Config.HTTPS {
		host := fmt.Sprintf("%s:%d", libs.Config.Host, 443)
		if err := s.App.Run(iris.TLS(host, libs.Config.Certpath, libs.Config.Certkey)); err != nil {
			return err
		}
	} else {
		if err := s.App.Run(
			iris.Addr(fmt.Sprintf("%s:%d", libs.Config.Host, libs.Config.Port)),
			iris.WithoutServerError(iris.ErrServerClosed),
			iris.WithOptimizations,
			iris.WithTimeFormat(time.RFC3339),
		); err != nil {
			return err
		}
	}
	s.Status = true
	return nil
}

// Start close the server at 3-6 seconds
func (s *Server) Stop() {
	go func() {
		time.Sleep(3 * time.Second)
		ctx, cancel := stdContext.WithTimeout(stdContext.TODO(), 3*time.Second)
		defer cancel()
		s.App.Shutdown(ctx)
		s.Status = false
	}()
}

// PathName
type PathName struct {
	Name   string
	Path   string
	Method string
}

// 获取路由信息
func (s *Server) GetRoutes() []*models.Permission {
	var rrs []*models.Permission
	names := getPathNames(s.App.GetRoutesReadOnly())
	if libs.Config.Debug {
		fmt.Println(fmt.Sprintf("路由权限集合：%v", names))
		fmt.Println(fmt.Sprintf("Iris App ：%v", s.App))
	}
	for _, pathName := range names {
		if !isPermRoute(pathName.Name) {
			rr := &models.Permission{Name: pathName.Path, DisplayName: pathName.Name, Description: pathName.Name, Act: pathName.Method}
			rrs = append(rrs, rr)
		}
	}
	return rrs
}

// getPathNames
func getPathNames(routeReadOnly []context.RouteReadOnly) []*PathName {
	var pns []*PathName
	if libs.Config.Debug {
		fmt.Println(fmt.Sprintf("routeReadOnly：%v", routeReadOnly))
	}
	for _, s := range routeReadOnly {
		pn := &PathName{
			Name:   s.Name(),
			Path:   s.Path(),
			Method: s.Method(),
		}
		pns = append(pns, pn)
	}

	return pns
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

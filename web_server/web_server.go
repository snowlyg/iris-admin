package web_server

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"net/http"
	"strings"
	"time"

	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/IrisAdminApi/config"
	"github.com/snowlyg/IrisAdminApi/models"
	"github.com/snowlyg/IrisAdminApi/routes"
	"github.com/snowlyg/IrisAdminApi/sysinit"
)

type Server struct {
	App       *iris.Application
	AssetFile http.FileSystem
}

func NewServer(assetFile http.FileSystem) *Server {
	app := iris.Default()
	return &Server{
		App:       app,
		AssetFile: assetFile,
	}
}

func (s *Server) Serve() error {
	if config.Config.HTTPS {
		host := fmt.Sprintf("%s:%d", config.Config.Host, 443)
		if err := s.App.Run(iris.TLS(host, config.Config.Certpath, config.Config.Certkey)); err != nil {
			return err
		}
	} else {
		if err := s.App.Run(
			iris.Addr(fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port)),
			iris.WithoutServerError(iris.ErrServerClosed),
			iris.WithOptimizations,
			iris.WithTimeFormat(time.RFC3339),
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) NewApp() {
	s.App.Logger().SetLevel(config.Config.LogLevel)

	// 二进制模式 ， 绑定前端文件
	if config.Config.Bindata {
		color.Yellow(`开启二进制模式前，请执行
		cd www
		npm run build:stage
		cd ..
		go generate 
		等步骤，打包前端文件
`)
		s.App.RegisterView(iris.HTML(s.AssetFile, ".html"))
		s.App.HandleDir("/", s.AssetFile)
	}

	db := sysinit.Db
	db.AutoMigrate(
		&models.User{},
		&models.OauthToken{},
		&models.Role{},
		&models.Permission{},
	)

	iris.RegisterOnInterrupt(func() {
		_ = db.Close()
	})

	routes.App(s.App) //注册 app 路由
}

type PathName struct {
	Name   string
	Path   string
	Method string
}

// 获取路由信息
func (s *Server) GetRoutes() []*models.Permission {
	var rrs []*models.Permission
	names := getPathNames(s.App.GetRoutesReadOnly())
	if config.Config.Debug {
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

func getPathNames(routeReadOnly []context.RouteReadOnly) []*PathName {
	var pns []*PathName
	if config.Config.Debug {
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

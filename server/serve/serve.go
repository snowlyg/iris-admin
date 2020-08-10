package serve

import (
	"fmt"
	"github.com/betacraft/yaag/yaag"
	"github.com/snowlyg/IrisAdminApi/backend/libs"
	"github.com/snowlyg/IrisAdminApi/backend/models"
	"github.com/snowlyg/IrisAdminApi/backend/routes"
	"github.com/snowlyg/IrisAdminApi/backend/sysinit"
	"path/filepath"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/backend/config"
)

type Server struct {
	App *iris.Application
}

func NewServer() *Server {
	app := iris.Default()
	return &Server{
		App: app,
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
	s.App.Logger().SetLevel("debug")

	s.App.RegisterView(iris.HTML(libs.WwwPath(), ".html"))

	db := sysinit.Db
	db.AutoMigrate(
		&models.User{},
		&models.OauthToken{},
		&models.Role{},
		&models.Permission{},
		&models.Stream{},
	)

	iris.RegisterOnInterrupt(func() {
		_ = db.Close()
	})

	docPath := filepath.Join(libs.WwwPath(), "apiDoc/index.html")
	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware. //api 文档配置
		On:       true,
		DocTitle: "GoIrisApi",
		DocPath:  docPath, //设置绝对路径
		BaseUrls: map[string]string{
			"Production": config.Config.Host,
			"Staging":    "",
		},
	})

	routes.App(s.App) //注册 app 路由
}

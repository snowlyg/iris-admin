package serve

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/betacraft/yaag/yaag"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/server/config"
	"github.com/snowlyg/IrisAdminApi/server/libs"
	"github.com/snowlyg/IrisAdminApi/server/models"
	"github.com/snowlyg/IrisAdminApi/server/routes"
	"github.com/snowlyg/IrisAdminApi/server/sysinit"
)

type Server struct {
	App        *iris.Application
	Asset      func(name string) ([]byte, error)
	AssetNames func() []string
	AssetInfo  func(name string) (os.FileInfo, error)
}

func NewServer(assetFn func(name string) ([]byte, error), namesFn func() []string, assetInfo func(name string) (os.FileInfo, error)) *Server {
	app := iris.Default()
	return &Server{
		App:        app,
		Asset:      assetFn,
		AssetNames: namesFn,
		AssetInfo:  assetInfo,
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

	tmpl := iris.HTML(libs.WwwPath(), ".html").Binary(s.Asset, s.AssetNames)
	s.App.RegisterView(tmpl)

	s.App.HandleDir("/", libs.WwwPath(), iris.DirOptions{
		Asset:      s.Asset,
		AssetInfo:  s.AssetInfo,
		AssetNames: s.AssetNames,
	})

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

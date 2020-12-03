package app

import (
	stdContext "context"
	"fmt"
	"github.com/snowlyg/blog/libs/logging"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/routes"
	"github.com/snowlyg/easygorm"
)

var Ser *Server

// Server
type Server struct {
	App    *iris.Application
	Status bool
}

func NewServer() *Server {
	app := iris.New()
	app.Logger().SetLevel(libs.Config.LogLevel)
	libs.InitRedisCluster(libs.GetRedisUris(), libs.Config.Redis.Pwd)
	routes.App(app)
	iris.RegisterOnInterrupt(func() {
		db, err := easygorm.Egm.Db.DB()
		if err != nil {
			logging.Err.Errorf("db init err: %+v", err)
			panic(err)
		}
		defer db.Close()
	})
	Ser = &Server{
		App:    app,
		Status: false,
	}
	return Ser
}

// Start
func (s *Server) Start() error {
	if libs.Config.HTTPS {
		host := fmt.Sprintf("%s:%d", libs.Config.Host, 443)
		if err := s.App.Run(iris.TLS(host, libs.Config.Certpath, libs.Config.Certkey)); err != nil {
			logging.Err.Errorf("app run https err: %+v", err)
			return err
		}
	} else {
		if err := s.App.Run(
			iris.Addr(fmt.Sprintf("%s:%d", libs.Config.Host, libs.Config.Port)),
			iris.WithoutServerError(iris.ErrServerClosed),
			iris.WithOptimizations,
			iris.WithTimeFormat(time.RFC3339),
		); err != nil {
			logging.Err.Errorf("app run http  err: %+v", err)
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

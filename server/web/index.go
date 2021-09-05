package web

import (
	stdContext "context"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/g"
	v1 "github.com/snowlyg/iris-admin/modules/v1"
	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/module"
	"github.com/snowlyg/iris-admin/server/viper"
	"github.com/snowlyg/iris-admin/server/zap"
)

type WebServer struct {
	app               *iris.Application  // iris application
	modules           []module.WebModule // 数据模型
	idleConnsClosed   chan struct{}
	addr              string //端口
	timeFormat        string // 时间格式
	globalMiddlewares []context.Handler
	wg                sync.WaitGroup
}

func Init() *WebServer {
	viper.Init()
	zap.Init()
	err := cache.Init()
	if err != nil {
		panic(err)
	}
	app := iris.New()
	app.Validator = validator.New() //参数验证
	app.Logger().SetLevel(g.CONFIG.System.Level)
	idleConnsClosed := make(chan struct{})
	iris.RegisterOnInterrupt(func() { //优雅退出
		timeout := 10 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		app.Shutdown(ctx) // close all hosts
		close(idleConnsClosed)
	})
	return &WebServer{
		app:               app,
		addr:              g.CONFIG.System.Addr,
		timeFormat:        g.CONFIG.System.TimeFormat,
		idleConnsClosed:   idleConnsClosed,
		globalMiddlewares: []context.Handler{},
	}
}

func (ws *WebServer) GetAddr() string {
	return ws.addr
}

func (ws *WebServer) AddModule(module ...module.WebModule) {
	ws.modules = append(ws.modules, module...)
}

func (ws *WebServer) AddStatic(requestPath string, fsOrDir interface{}, opts ...router.DirOptions) {
	ws.app.HandleDir(requestPath, fsOrDir, opts...)
}

func (ws *WebServer) GetModules() []module.WebModule {
	return ws.modules
}

func (ws *WebServer) Run() {
	if ws.addr == "" { // 默认 8085
		ws.addr = "127.0.0.1:8085"
	}
	if ws.timeFormat == "" { // 默认 80
		ws.timeFormat = time.RFC3339
	}
	ws.app.UseGlobal(ws.globalMiddlewares...)
	ws.AddModule(v1.Party())
	ws.app.HandleDir("/static", iris.Dir(filepath.Join(dir.GetCurrentAbPath(), "static")))
	ws.InitRouter()
	ws.app.Listen(
		ws.addr,
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(ws.timeFormat),
	)
	<-ws.idleConnsClosed
}

package web

import (
	stdContext "context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/helper/tests"
	"github.com/snowlyg/iris-admin/g"
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
	if g.CONFIG.System.Addr == "" { // 默认 8085
		g.CONFIG.System.Addr = "127.0.0.1:8085"
	}
	if g.CONFIG.System.TimeFormat == "" { // 默认 80
		g.CONFIG.System.TimeFormat = time.RFC3339
	}
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

func (ws *WebServer) AddStatic(requestPath string, fsOrDir interface{}, opts ...iris.DirOptions) {
	ws.app.HandleDir(requestPath, fsOrDir, opts...)
}

func (ws *WebServer) GetModules() []module.WebModule {
	return ws.modules
}

func (ws *WebServer) GetTestAuth(t *testing.T) *tests.Client {
	client := tests.New(str.Join("http://", ws.addr), t, ws.app)
	if client == nil {
		t.Fatalf("client is nil")
	}
	return client
}

func (ws *WebServer) GetTestLogin(t *testing.T, url string, res tests.Responses, datas ...map[string]interface{}) *httpexpect.Expect {
	return ws.GetTestAuth(t).Login(url, res, datas...)
}

func (ws *WebServer) GetTestLogout(t *testing.T, url string, res tests.Responses) {
	ws.GetTestAuth(t).Logout(url, res)
}

func (ws *WebServer) Run() {
	ws.app.UseGlobal(ws.globalMiddlewares...)
	err := ws.InitRouter()
	if err != nil {
		fmt.Printf("初始化路由错误： %v\n", err)
		panic(err)
	}
	ws.app.Listen(
		ws.addr,
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(ws.timeFormat),
	)
	<-ws.idleConnsClosed
}

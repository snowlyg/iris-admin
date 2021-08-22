package web

import (
	stdContext "context"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/iris-admin/g"
)

type WebServer struct {
	app               *iris.Application // iris application
	modules           []WebModule       // 数据模型
	idleConnsClosed   chan struct{}
	addr              string //端口
	timeFormat        string // 时间格式
	globalMiddlewares []context.Handler
}

func Init() *WebServer {
	app := iris.New()
	app.Logger().SetLevel(g.CONFIG.System.Level)
	idleConnsClosed := make(chan struct{})
	iris.RegisterOnInterrupt(func() { //优雅退出
		timeout := 10 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		app.Shutdown(ctx) // close all hosts
		close(idleConnsClosed)
	})
	return &WebServer{app: app, addr: g.CONFIG.System.Addr, timeFormat: g.CONFIG.System.TimeFormat, idleConnsClosed: idleConnsClosed}
}

func (ws *WebServer) GetAddr() string {
	return ws.addr
}

func (ws *WebServer) AddModule(module WebModule) {
	ws.modules = append(ws.modules, module)
}

func (ws *WebServer) GetModules() []WebModule {
	return ws.modules
}

func (ws *WebServer) Run() {
	if ws.addr == "" { // 默认 8085
		ws.addr = "127.0.0.1:8085"
	}
	if ws.timeFormat == "" { // 默认 80
		ws.timeFormat = time.RFC3339
	}
	fmt.Printf("listen on %s", ws.addr)
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

package web

import (
	stdContext "context"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/helper/dir"
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
	staticPrefix      string //静态文件访问地址前缀
	staticPath        string //静态文件地址
	webPath           string //前端文件地址
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
	if g.CONFIG.System.StaticPath == "" { // 默认 /static/upload
		g.CONFIG.System.StaticPath = "/static/upload"
	}
	if g.CONFIG.System.StaticPrefix == "" { // 默认 /upload
		g.CONFIG.System.StaticPrefix = "/upload"
	}
	if g.CONFIG.System.WebPath == "" { // 默认 /./dist
		g.CONFIG.System.WebPath = "./dist"
	}
	if g.CONFIG.System.TimeFormat == "" { // 默认 80
		g.CONFIG.System.TimeFormat = time.RFC3339
	}
	return &WebServer{
		app:               app,
		addr:              g.CONFIG.System.Addr,
		timeFormat:        g.CONFIG.System.TimeFormat,
		staticPrefix:      g.CONFIG.System.StaticPrefix,
		staticPath:        g.CONFIG.System.StaticPath,
		webPath:           g.CONFIG.System.WebPath,
		idleConnsClosed:   idleConnsClosed,
		globalMiddlewares: []context.Handler{},
	}
}

func (ws *WebServer) GetStaticPath() string {
	return ws.staticPath
}

func (ws *WebServer) GetWebPath() string {
	return ws.webPath
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

func (ws *WebServer) AddWebStatic(requestPath string) {
	fsOrDir := iris.Dir(filepath.Join(dir.GetCurrentAbPath(), ws.webPath))
	ws.AddStatic(requestPath, fsOrDir, iris.DirOptions{
		IndexName: "index.html",
		SPA:       true,
	})
}

func (ws *WebServer) AddUploadStatic() {
	fsOrDir := iris.Dir(filepath.Join(dir.GetCurrentAbPath(), ws.staticPath))
	ws.AddStatic(ws.staticPrefix, fsOrDir)
}

func (ws *WebServer) GetModules() []module.WebModule {
	return ws.modules
}

var client *tests.Client

func (ws *WebServer) GetTestAuth(t *testing.T) *tests.Client {
	var once sync.Once
	once.Do(
		func() {
			client = tests.New(str.Join("http://", ws.addr), t, ws.app)
			if client == nil {
				t.Fatalf("client is nil")
			}
		},
	)

	return client
}

func (ws *WebServer) GetTestLogin(t *testing.T, url string, res tests.Responses, datas ...map[string]interface{}) *tests.Client {
	client := ws.GetTestAuth(t)
	err := client.Login(url, res, datas...)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func (ws *WebServer) Run() {
	ws.app.UseGlobal(ws.globalMiddlewares...)
	err := ws.InitRouter()
	if err != nil {
		fmt.Printf("初始化路由错误： %v\n", err)
		panic(err)
	}
	// 添加上传文件路径
	ws.app.Listen(
		ws.addr,
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(ws.timeFormat),
	)
	<-ws.idleConnsClosed
}

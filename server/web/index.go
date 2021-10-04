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
	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/module"
	"github.com/snowlyg/iris-admin/server/viper_server"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"github.com/snowlyg/multi"
)

var client *tests.Client

// WebServer web 服务
// - app iris application
// - modules 服务的模块
// - idleConnsClosed
// - addr  服务访问地址
// - timeFormat  时间格式
// - globalMiddlewares  全局中间件
// - wg  sync.WaitGroup
// - staticPrefix  静态文件访问地址前缀
// - staticPath  静态文件地址
// - webPath  前端文件地址
type WebServer struct {
	app               *iris.Application
	modules           []module.WebModule
	idleConnsClosed   chan struct{}
	addr              string
	timeFormat        string
	globalMiddlewares []context.Handler
	wg                sync.WaitGroup
	staticPrefix      string
	staticPath        string
	webPath           string
}

// Init 初始化web服务
func Init() *WebServer {
	viper_server.Init(getViperConfig())
	zap_server.Init()

	// 初始化认证
	err := multi.InitDriver(
		&multi.Config{
			DriverType:      CONFIG.System.CacheType,
			UniversalClient: cache.Instance()},
	)
	if err != nil || multi.AuthDriver == nil {
		panic(fmt.Sprintf("认证驱动初始化失败 %v \n", err))
	}

	app := iris.New()
	app.Validator = validator.New() //参数验证
	app.Logger().SetLevel(CONFIG.System.Level)
	idleConnsClosed := make(chan struct{})
	iris.RegisterOnInterrupt(func() { //优雅退出
		timeout := 10 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		app.Shutdown(ctx) // close all hosts
		close(idleConnsClosed)
	})

	if CONFIG.System.Addr == "" { // 默认 8085
		CONFIG.System.Addr = "127.0.0.1:8085"
	}

	if CONFIG.System.StaticPath == "" { // 默认 /static/upload
		CONFIG.System.StaticPath = "/static/upload"
	}

	if CONFIG.System.StaticPrefix == "" { // 默认 /upload
		CONFIG.System.StaticPrefix = "/upload"
	}

	if CONFIG.System.WebPath == "" { // 默认 /./dist
		CONFIG.System.WebPath = "./dist"
	}

	if CONFIG.System.TimeFormat == "" { // 默认 80
		CONFIG.System.TimeFormat = time.RFC3339
	}

	return &WebServer{
		app:               app,
		addr:              CONFIG.System.Addr,
		timeFormat:        CONFIG.System.TimeFormat,
		staticPrefix:      CONFIG.System.StaticPrefix,
		staticPath:        CONFIG.System.StaticPath,
		webPath:           CONFIG.System.WebPath,
		idleConnsClosed:   idleConnsClosed,
		globalMiddlewares: []context.Handler{},
	}
}

// GetStaticPath 获取静态路径
func (ws *WebServer) GetStaticPath() string {
	return ws.staticPath
}

// GetWebPath 获取前端路径
func (ws *WebServer) GetWebPath() string {
	return ws.webPath
}

// GetAddr 获取web服务地址
func (ws *WebServer) GetAddr() string {
	return ws.addr
}

// AddModule 添加模块
func (ws *WebServer) AddModule(module ...module.WebModule) {
	ws.modules = append(ws.modules, module...)
}

// AddStatic 添加静态文件
func (ws *WebServer) AddStatic(requestPath string, fsOrDir interface{}, opts ...iris.DirOptions) {
	ws.app.HandleDir(requestPath, fsOrDir, opts...)
}

// AddWebStatic 添加前端访问地址
func (ws *WebServer) AddWebStatic(requestPath string) {
	fsOrDir := iris.Dir(filepath.Join(dir.GetCurrentAbPath(), ws.webPath))
	ws.AddStatic(requestPath, fsOrDir, iris.DirOptions{
		IndexName: "index.html",
		SPA:       true,
	})
}

// AddUploadStatic 添加上传文件访问地址
func (ws *WebServer) AddUploadStatic() {
	fsOrDir := iris.Dir(filepath.Join(dir.GetCurrentAbPath(), ws.staticPath))
	ws.AddStatic(ws.staticPrefix, fsOrDir)
}

// GetModules 获取模块
func (ws *WebServer) GetModules() []module.WebModule {
	return ws.modules
}

// GetTestAuth 获取测试验证客户端
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

// GetTestLogin 测试登录web服务
func (ws *WebServer) GetTestLogin(t *testing.T, url string, res tests.Responses, datas ...map[string]interface{}) *tests.Client {
	client := ws.GetTestAuth(t)
	err := client.Login(url, res, datas...)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// Run 启动web服务
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

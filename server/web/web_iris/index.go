package web_iris

import (
	stdContext "context"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/httptest"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_iris/middleware"
)

var ErrAuthDriverEmpty = errors.New("认证驱动初始化失败")

// WebServer web服务
// - app iris application
// - idleConnsClosed
// - addr  服务访问地址
// - timeFormat  时间格式
// - staticPrefix  静态文件访问地址前缀
// - staticPath  静态文件地址
// - webPath  前端文件地址
type WebServer struct {
	app             *iris.Application
	idleConnsClosed chan struct{}
	parties         []Party
	addr            string
	timeFormat      string
	staticPrefix    string
	staticPath      string
	webPrefix       string
	webPath         string
}

// Party 功能模块
// - perfix 模块路由路径
// - partyFunc 模块
type Party struct {
	Perfix    string
	PartyFunc func(index iris.Party)
}

// Init 初始化web服务
// 先初始化基础服务 config , zap , database , casbin  e.g.
func Init() *WebServer {
	web.InitWeb()
	app := iris.New()
	if web.CONFIG.System.Tls {
		app.Use(middleware.LoadTls()) // 打开就能玩https了
	}
	app.Use(recover.New())
	app.Validator = validator.New() //参数验证
	app.Logger().SetLevel(web.CONFIG.System.Level)
	idleConnsClosed := make(chan struct{})
	iris.RegisterOnInterrupt(func() { //优雅退出
		timeout := 10 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		app.Shutdown(ctx) // close all hosts
		close(idleConnsClosed)
	})

	web.Verfiy()

	return &WebServer{
		app:             app,
		addr:            web.CONFIG.System.Addr,
		timeFormat:      web.CONFIG.System.TimeFormat,
		staticPrefix:    web.CONFIG.System.StaticPrefix,
		staticPath:      web.CONFIG.System.StaticPath,
		webPrefix:       web.CONFIG.System.WebPrefix,
		webPath:         web.CONFIG.System.WebPath,
		idleConnsClosed: idleConnsClosed,
	}
}

// AddStatic 添加静态文件
func (ws *WebServer) AddStatic(requestPath string, fsOrDir interface{}, opts ...iris.DirOptions) {
	ws.app.HandleDir(requestPath, fsOrDir, opts...)
}

// AddModule 添加模块
func (ws *WebServer) AddModule(parties ...Party) {
	ws.parties = append(ws.parties, parties...)
}

// AddWebStatic 添加前端访问地址
func (ws *WebServer) AddWebStatic() {
	fsOrDir := iris.Dir(filepath.Join(dir.GetCurrentAbPath(), ws.webPath))
	opt := iris.DirOptions{
		IndexName: "index.html",
		SPA:       true,
	}
	ws.AddStatic(ws.webPrefix, fsOrDir, opt)
}

// AddUploadStatic 添加上传文件访问地址
func (ws *WebServer) AddUploadStatic() {
	fsOrDir := iris.Dir(filepath.Join(dir.GetCurrentAbPath(), ws.staticPath))
	ws.AddStatic(ws.staticPrefix, fsOrDir)
}

// GetTestClient 获取测试验证客户端
func (ws *WebServer) GetTestClient(t *testing.T) *httptest.Client {
	if ws.app == nil {
		t.Errorf("ws.app is nil")
	}
	var once sync.Once
	var client *httptest.Client
	once.Do(
		func() {
			client = httptest.New(str.Join("http://", ws.addr), t, ws.app)
			if client == nil {
				t.Errorf("test client is nil")
			}
		},
	)

	return client
}

// GetTestLogin 测试登录web服务
func (ws *WebServer) GetTestLogin(t *testing.T, url string, res httptest.Responses, datas ...interface{}) *httptest.Client {
	client := ws.GetTestClient(t)
	if client == nil {
		t.Error("登录失败")
		return nil
	}
	err := client.Login(url, res, datas...)
	if err != nil {
		t.Error(err)
		return nil
	}
	return client
}

// Run 启动web服务
func (ws *WebServer) Run() {
	ws.app.Listen(
		ws.addr,
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(ws.timeFormat),
	)
	<-ws.idleConnsClosed
}

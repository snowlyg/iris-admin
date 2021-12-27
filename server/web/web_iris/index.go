package web_iris

import (
	stdContext "context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/httptest"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_iris/middleware"
	"github.com/snowlyg/iris-admin/server/zap_server"
)

var ErrAuthDriverEmpty = errors.New("认证驱动初始化失败")

// WebServer web服务
// - app iris application
// - idleConnsClosed
// - addr  服务访问地址
// - timeFormat  时间格式
// - staticPrefix  静态文件访问地址前缀

type WebServer struct {
	app             *iris.Application
	idleConnsClosed chan struct{}
	parties         []Party
	addr            string
	timeFormat      string
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
		idleConnsClosed: idleConnsClosed,
	}
}

// AddModule 添加模块
func (ws *WebServer) AddModule(parties ...Party) {
	ws.parties = append(ws.parties, parties...)
}

// AddWebStatic 添加前端访问地址
func (ws *WebServer) AddWebStatic(paths ...string) {
	if len(paths) != 2 {
		zap_server.ZAPLOG.Warn("AddWebStatic function need 2 params")
		return
	}

	if paths[0] == "" || paths[1] == "" {
		zap_server.ZAPLOG.Warn("AddWebStatic function params not support empty string")
		return
	}
	webPrefix := paths[0]
	webPrefixs := strings.Split(web.CONFIG.System.WebPrefix, ",")
	if str.InStrArray(webPrefix, webPrefixs) {
		return
	}
	staticAbsPath := paths[1]
	fsOrDir := iris.Dir(staticAbsPath)
	opt := iris.DirOptions{
		IndexName: "index.html",
		SPA:       true,
	}
	ws.app.HandleDir(webPrefix, fsOrDir, opt)
	web.CONFIG.System.WebPrefix = str.Join(web.CONFIG.System.WebPrefix, ",", webPrefix)
}

// AddUploadStatic 添加上传文件访问地址
func (ws *WebServer) AddUploadStatic(paths ...string) {
	if len(paths) != 2 {
		zap_server.ZAPLOG.Warn("AddWebStatic function need 2 params")
		return
	}

	if paths[0] == "" || paths[1] == "" {
		zap_server.ZAPLOG.Warn("AddWebStatic function params not support empty string")
		return
	}
	staticPrefix := paths[0]
	staticAbsPath := paths[1]
	fsOrDir := iris.Dir(staticAbsPath)
	ws.app.HandleDir(staticPrefix, fsOrDir)
	web.CONFIG.System.StaticPrefix = staticPrefix
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

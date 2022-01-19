package web_iris

import (
	stdContext "context"
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/snowlyg/helper/str"
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

// GetEngine 增加灵活性
func (ws *WebServer) GetEngine() *iris.Application {
	return ws.app
}

// AddModule 添加模块
func (ws *WebServer) AddModule(parties ...Party) {
	ws.parties = append(ws.parties, parties...)
}

// AddWebStatic 添加前端访问地址
func (ws *WebServer) AddWebStatic(staticAbsPath, webPrefix string, paths ...string) {
	webPrefixs := strings.Split(web.CONFIG.System.WebPrefix, ",")
	if str.InStrArray(webPrefix, webPrefixs) {
		return
	}

	fsOrDir := iris.Dir(staticAbsPath)
	opt := iris.DirOptions{
		IndexName: "index.html",
		SPA:       true,
	}
	ws.app.HandleDir(webPrefix, fsOrDir, opt)
	web.CONFIG.System.WebPrefix = str.Join(web.CONFIG.System.WebPrefix, ",", webPrefix)
}

// AddUploadStatic 添加上传文件访问地址
func (ws *WebServer) AddUploadStatic(webPrefix, staticAbsPath string) {
	fsOrDir := iris.Dir(staticAbsPath)
	ws.app.HandleDir(webPrefix, fsOrDir)
	web.CONFIG.System.StaticPrefix = webPrefix
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

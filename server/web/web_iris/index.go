package web_iris

import (
	stdContext "context"
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/snowlyg/helper/arr"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_iris/middleware"
)

var ErrAuthDriverEmpty = errors.New("auth driver initialize fail")

// WebServer
// - app iris application
// - idleConnsClosed
// - addr
// - timeFormat
// - staticPrefix

type WebServer struct {
	app             *iris.Application
	idleConnsClosed chan struct{}
	parties         []Party
	addr            string
	timeFormat      string
}

// Party
// - perfix
// - partyFunc
type Party struct {
	Perfix    string
	PartyFunc func(index iris.Party)
}

// Init
func Init() *WebServer {
	app := iris.New()
	if web.CONFIG.System.Tls {
		app.Use(middleware.LoadTls())
	}
	app.Use(recover.New())
	app.Validator = validator.New()
	app.Logger().SetLevel(web.CONFIG.System.Level)
	idleConnsClosed := make(chan struct{})
	iris.RegisterOnInterrupt(func() {
		timeout := 10 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		app.Shutdown(ctx) // close all hosts
		close(idleConnsClosed)
	})

	web.SetDefaultAddrAndTimeFormat()

	return &WebServer{
		app:             app,
		addr:            web.CONFIG.System.Addr,
		timeFormat:      web.CONFIG.System.TimeFormat,
		idleConnsClosed: idleConnsClosed,
	}
}

// GetEngine
func (ws *WebServer) GetEngine() *iris.Application {
	return ws.app
}

// AddModule
func (ws *WebServer) AddModule(parties ...Party) {
	ws.parties = append(ws.parties, parties...)
}

// AddWebStatic
func (ws *WebServer) AddWebStatic(staticAbsPath, webPrefix string, paths ...string) {
	webPrefixs := strings.Split(web.CONFIG.System.WebPrefix, ",")
	wp := arr.NewCheckArrayType(2)
	for _, webPrefix := range webPrefixs {
		wp.Add(webPrefix)
	}
	if wp.Check(webPrefix) {
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

// AddUploadStatic
func (ws *WebServer) AddUploadStatic(webPrefix, staticAbsPath string) {
	fsOrDir := iris.Dir(staticAbsPath)
	ws.app.HandleDir(webPrefix, fsOrDir)
	web.CONFIG.System.StaticPrefix = webPrefix
}

// Run
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

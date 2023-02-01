package web_gin

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"github.com/snowlyg/helper/arr"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_gin/middleware"
)

var ErrAuthDriverEmpty = errors.New("auth driver initialize fail")

// WebServer
// - app gin.Engine
// - idleConnsClosed
// - addr
// - timeFormat
// - staticPrefix
type WebServer struct {
	app *gin.Engine
	server
	addr       string
	timeFormat string
	webStatics []WebStatic
}

type WebStatic struct {
	Prefix    string
	IndexFile []byte
}

// Init
func Init() *WebServer {
	gin.SetMode(web.CONFIG.System.Level)
	app := gin.Default()
	if web.CONFIG.System.Tls {
		app.Use(middleware.LoadTls())
	}
	app.Use(middleware.Cors())
	registerValidation()

	gin.DefaultWriter = colorable.NewColorableStdout()

	web.Verfiy()

	return &WebServer{
		app:        app,
		addr:       web.CONFIG.System.Addr,
		timeFormat: web.CONFIG.System.TimeFormat,
	}
}

// NoRoute for 404 http status
func (ws *WebServer) NoRoute() {
	if len(ws.webStatics) == 0 {
		return
	}

	ws.app.NoRoute(func(ctx *gin.Context) {
		// excepte for /v0 /v1 and so on
		reg := `^/v[0-9]+$|^(/v[0-9]+)/`
		ok, _ := regexp.MatchString(reg, ctx.Request.RequestURI)
		if ok {
			ctx.Writer.WriteHeader(http.StatusNotFound)
			ctx.Writer.Flush()
			return
		}

		var indexFile []byte
		for _, wp := range ws.webStatics {
			// match /admin or /admin/***
			reg := str.Join("^", wp.Prefix, "$|^(", wp.Prefix, ")/")
			ok, err := regexp.MatchString(reg, ctx.Request.RequestURI)
			if err != nil || !ok {
				continue
			}
			indexFile = wp.IndexFile
		}

		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write(indexFile)

		ctx.Writer.Header().Add("Accept", "text/html")
		ctx.Writer.Flush()
	})
}

// GetEngine return *gin.Engine
func (ws *WebServer) GetEngine() *gin.Engine {
	return ws.app
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

	favicon := filepath.Join(staticAbsPath, "favicon.ico")
	index := filepath.Join(staticAbsPath, "index.html")

	ws.app.Static(str.Join(webPrefix, "/favicon.ico"), favicon)
	ws.app.StaticFile(webPrefix, index)

	if len(paths) > 0 {
		for _, path := range paths {
			static := filepath.Join(staticAbsPath, path)
			ws.app.Static(path, static)
		}
	}

	web.CONFIG.System.WebPrefix = str.Join(web.CONFIG.System.WebPrefix, ",", webPrefix)
	file, _ := dir.ReadBytes(index)
	webStatic := WebStatic{
		Prefix:    webPrefix,
		IndexFile: file,
	}
	ws.webStatics = append(ws.webStatics, webStatic)

}

// AddUploadStatic
func (ws *WebServer) AddUploadStatic(webPrefix, staticAbsPath string) {
	ws.app.StaticFS(webPrefix, http.Dir(staticAbsPath))
	web.CONFIG.System.StaticPrefix = webPrefix
}

// Run
func (ws *WebServer) Run() {
	ws.NoRoute()
	s := initServer(web.CONFIG.System.Addr, ws.app)
	time.Sleep(10 * time.Microsecond)
	fmt.Printf("默认监听地址: http://%s\n", web.CONFIG.System.Addr)
	s.ListenAndServe()

}

package admin

import (
	"errors"
	"fmt"
	"log"
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
	"github.com/snowlyg/iris-admin/conf"
)

// Start
func Start(wf *WebServe) error {
	if err := wf.InitRouter(); err != nil {
		return fmt.Errorf("init router fail:%s", err.Error())
	}
	wf.Run()
	return nil
}

// StartTest
func StartTest(wf *WebServe) {
	err := wf.InitRouter()
	if err != nil {
		log.Printf("start test fail:%s\n", err.Error())
	}
}

var ErrAuthDriverEmpty = errors.New("auth driver initialize fail")

// WebServe
// - app gin.Engine
// - idleConnsClosed
// - addr
// - timeFormat
// - staticPrefix
type WebServe struct {
	server
	app  *gin.Engine
	conf *conf.Conf

	addr       string
	timeFormat string

	statics []WebStatic
}

type WebStatic struct {
	Prefix    string
	IndexFile []byte
}

// NewServe
func NewServe() *WebServe {
	config := conf.NewConf()
	gin.SetMode(config.System.GinMode)
	app := gin.Default()
	if config.System.Tls {
		app.Use(LoadTls())
	}
	app.Use(config.CorsConf.Cors())
	registerValidation()

	gin.DefaultWriter = colorable.NewColorableStdout()

	config.SetDefaultAddrAndTimeFormat()

	return &WebServe{
		app:        app,
		addr:       config.System.Addr,
		timeFormat: config.System.TimeFormat,
	}
}

// NoRoute for 404 http status
func (ws *WebServe) NoRoute() {
	if len(ws.statics) == 0 {
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
		for _, wp := range ws.statics {
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
func (ws *WebServe) GetEngine() *gin.Engine {
	return ws.app
}

// AddWebStatic
func (ws *WebServe) AddWebStatic(staticAbsPath, webPrefix string, paths ...string) {
	webPrefixs := strings.Split(ws.conf.System.WebPrefix, ",")
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

	ws.conf.System.WebPrefix = str.Join(ws.conf.System.WebPrefix, ",", webPrefix)
	file, _ := dir.ReadBytes(index)
	webStatic := WebStatic{
		Prefix:    webPrefix,
		IndexFile: file,
	}
	ws.statics = append(ws.statics, webStatic)

}

// AddUploadStatic
func (ws *WebServe) AddUploadStatic(webPrefix, staticAbsPath string) {
	ws.app.StaticFS(webPrefix, http.Dir(staticAbsPath))
	ws.conf.System.StaticPrefix = webPrefix
}

// Run
func (ws *WebServe) Run() {
	ws.NoRoute()
	s := initServer(ws.conf.System.Addr, ws.app)
	time.Sleep(10 * time.Microsecond)
	fmt.Printf("默认监听地址: http://%s\n", ws.conf.System.Addr)
	s.ListenAndServe()

}

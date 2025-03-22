package admin

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"github.com/snowlyg/helper/arr"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/conf"
	"gorm.io/gorm"
)

// // Start
// func Start(wf *WebServe) error {
// 	if err := wf.InitRouter(); err != nil {
// 		return fmt.Errorf("init router fail:%s", err.Error())
// 	}
// 	wf.Run()
// 	return nil
// }

// // StartTest
// func StartTest(wf *WebServe) {
// 	err := wf.InitRouter()
// 	if err != nil {
// 		log.Printf("start test fail:%s\n", err.Error())
// 	}
// }

// WebServe
// - app gin.Engine
// - idleConnsClosed
// - addr
// - timeFormat
// - staticPrefix
type WebServe struct {
	serve
	db       *gorm.DB
	enforcer *casbin.Enforcer
	engine   *gin.Engine
	conf     *conf.Conf

	statics []WebStatic
	// addr       string
	// timeFormat string
}

// WebStatic
type WebStatic struct {
	Prefix    string
	IndexFile []byte
}

// NewServe
func NewServe() (*WebServe, error) {
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
	db, err := gormDb(&config.Mysql)
	if err != nil {
		return nil, err
	}
	auth, err := getEnforcer(db)
	if err != nil {
		return nil, err
	}
	return &WebServe{
		conf:     config,
		engine:   app,
		enforcer: auth,
		db:       db,
	}, nil
}

// Engine return *gin.Engine
func (ws *WebServe) Engine() *gin.Engine {
	return ws.engine
}

// Config
func (ws *WebServe) Config() *conf.Conf {
	return ws.conf
}

// Auth
func (ws *WebServe) Auth() *casbin.Enforcer {
	return ws.enforcer
}

// DB
func (ws *WebServe) DB() *gorm.DB {
	return ws.db
}

// AddWebStatic
func (ws *WebServe) AddWebStatic(staticAbsPath, webPrefix string, paths ...string) {
	if ws.Engine() == nil {
		return
	}
	webPrefixs := strings.Split(ws.Config().System.WebPrefix, ",")
	wp := arr.NewCheckArrayType(2)
	for _, webPrefix := range webPrefixs {
		wp.Add(webPrefix)
	}
	if wp.Check(webPrefix) {
		return
	}

	favicon := filepath.Join(staticAbsPath, "favicon.ico")
	index := filepath.Join(staticAbsPath, "index.html")

	ws.Engine().Static(str.Join(webPrefix, "/favicon.ico"), favicon)
	ws.Engine().StaticFile(webPrefix, index)

	if len(paths) > 0 {
		for _, path := range paths {
			static := filepath.Join(staticAbsPath, path)
			ws.Engine().Static(path, static)
		}
	}

	ws.Config().System.WebPrefix = str.Join(ws.Config().System.WebPrefix, ",", webPrefix)
	file, _ := dir.ReadBytes(index)
	webStatic := WebStatic{
		Prefix:    webPrefix,
		IndexFile: file,
	}
	ws.statics = append(ws.statics, webStatic)

}

// AddUploadStatic
func (ws *WebServe) AddUploadStatic(webPrefix, staticAbsPath string) {
	if ws.Engine() != nil {
		ws.Engine().StaticFS(webPrefix, http.Dir(staticAbsPath))
	}
	if ws.Config() != nil {
		ws.Config().System.StaticPrefix = webPrefix
	}
}

// Run
func (ws *WebServe) Run() {
	if len(ws.statics) == 0 {
		return
	}
	if ws.Engine() == nil {
		return
	}

	ws.Engine().NoRoute(func(ctx *gin.Context) {
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
	s := run(ws.Config().System.Addr, ws.engine)
	time.Sleep(10 * time.Microsecond)
	fmt.Printf("默认监听地址: http://%s\n", ws.Config().System.Addr)
	s.ListenAndServe()
}

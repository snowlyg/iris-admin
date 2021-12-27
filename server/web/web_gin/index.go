package web_gin

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/httptest"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_gin/middleware"
)

var ErrAuthDriverEmpty = errors.New("认证驱动初始化失败")

// WebServer web服务
// - app gin.Engine
// - idleConnsClosed
// - addr  服务访问地址
// - timeFormat  时间格式
// - staticPrefix  静态文件访问地址前缀
// - staticAbsPath  静态文件地址
type WebServer struct {
	app *gin.Engine
	server
	addr          string
	timeFormat    string
	staticPrefix  string
	staticAbsPath string
	webPrefix     string
}

// Init 初始化web服务
// 先初始化基础服务 config , zap , database , casbin  e.g.
func Init() *WebServer {
	web.InitWeb()
	gin.SetMode(web.CONFIG.System.Level)
	app := gin.Default()
	if web.CONFIG.System.Tls {
		app.Use(middleware.LoadTls()) // 打开就能玩https了
	}
	registerValidation()

	web.Verfiy()

	return &WebServer{
		app:           app,
		addr:          web.CONFIG.System.Addr,
		timeFormat:    web.CONFIG.System.TimeFormat,
		staticPrefix:  web.CONFIG.System.StaticPrefix,
		staticAbsPath: web.CONFIG.System.StaticAbsPath,
		webPrefix:     web.CONFIG.System.WebPrefix,
	}
}

// AddWebStatic 添加前端访问地址
func (ws *WebServer) AddWebStatic() {
	favicon := filepath.Join(ws.staticAbsPath, "favicon.ico")
	index := filepath.Join(ws.staticAbsPath, "index.html")
	ws.app.Static("/favicon.ico", favicon)
	filepathNames, _ := filepath.Glob(filepath.Join(ws.staticAbsPath, "*"))
	for _, filepathName := range filepathNames {
		if filepathName == ws.staticAbsPath {
			continue
		}
		if dir.IsFile(filepathName) {
			continue
		}
		ws.app.Static(filepath.Base(filepathName), filepathName)
	}

	// 关键点【解决页面刷新404的问题】
	ws.app.NoRoute(func(ctx *gin.Context) {
		ctx.Writer.WriteHeader(http.StatusOK)
		if strings.Contains(ctx.Request.RequestURI, ws.webPrefix) {
			file, _ := dir.ReadBytes(index)
			ctx.Writer.Write(file)
		}
		ctx.Writer.Header().Add("Accept", "text/html")
		ctx.Writer.Flush()
	})
}

// AddUploadStatic 添加上传文件访问地址
func (ws *WebServer) AddUploadStatic() {
	ws.app.StaticFS(ws.staticPrefix, http.Dir(ws.staticPrefix))
}

// GetTestClient 获取测试验证客户端
func (ws *WebServer) GetTestClient(t *testing.T) *httptest.Client {
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
	s := initServer(web.CONFIG.System.Addr, ws.app)
	time.Sleep(10 * time.Microsecond)
	fmt.Printf("默认监听地址:http://%s\n", web.CONFIG.System.Addr)
	s.ListenAndServe()

}

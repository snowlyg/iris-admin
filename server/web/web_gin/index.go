package web_gin

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/httptest"
	"github.com/snowlyg/iris-admin/migration"
	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_gin/middleware"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"github.com/snowlyg/multi"
	"go.uber.org/zap"
)

var ErrAuthDriverEmpty = errors.New("认证驱动初始化失败")

// WebServer web服务
// - app gin.Engine
// - idleConnsClosed
// - addr  服务访问地址
// - timeFormat  时间格式
// - staticPrefix  静态文件访问地址前缀
// - staticPath  静态文件地址
// - webPath  前端文件地址
type WebServer struct {
	app *gin.Engine
	server
	addr         string
	timeFormat   string
	staticPrefix string
	staticPath   string
	webPrefix    string
	webPath      string
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
		app:          app,
		addr:         web.CONFIG.System.Addr,
		timeFormat:   web.CONFIG.System.TimeFormat,
		staticPrefix: web.CONFIG.System.StaticPrefix,

		webPrefix: web.CONFIG.System.WebPrefix,
		webPath:   web.CONFIG.System.WebPath,
	}
}

// AddStatic 添加静态文件
func (ws *WebServer) AddStatic(requestPath, root string) {
	ws.app.Static(requestPath, root)
}

// AddWebStatic 添加前端访问地址
func (ws *WebServer) AddWebStatic() {
	ws.AddStatic(ws.webPrefix, ws.webPath)
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

func BeforeTestMain(mysqlPwd, redisPwd string, redisDB int, party func(wi *WebServer), seed func(wi *WebServer, mc *migration.MigrationCmd)) (string, *WebServer) {
	fmt.Println("+++++ before test +++++")
	node, _ := snowflake.NewNode(1)
	uuid := str.Join("gin", "_", node.Generate().String())
	fmt.Printf("+++++ %s +++++\n\n", uuid)
	web.CONFIG.System.DbType = "mysql"
	web.InitWeb()

	database.CONFIG.Dbname = uuid
	database.CONFIG.Password = strings.TrimSpace(mysqlPwd)
	database.CONFIG.LogMode = true
	database.InitMysql()

	cache.CONFIG.DB = redisDB
	cache.CONFIG.Password = strings.TrimSpace(redisPwd)
	cache.InitCache()

	wi := Init()
	party(wi)
	web.StartTest(wi)

	mc := migration.New()
	// 添加 v1 内置模块数据表和数据
	fmt.Println("++++++ add model ++++++")
	seed(wi, mc)
	err := mc.Migrate()
	if err != nil {
		fmt.Printf("migrate get error [%s]", err.Error())
		return uuid, nil
	}
	err = mc.Seed()
	if err != nil {
		fmt.Printf("seed get error [%s]", err.Error())
		return uuid, nil
	}

	return uuid, wi
}

func AfterTestMain(uuid, logOutUrl string, testClient *httptest.Client) {
	fmt.Println("++++++++ after test main ++++++++")
	err := database.DorpDB(database.CONFIG.BaseDsn(), "mysql", uuid)
	if err != nil {
		text := str.Join("删除数据库 '", uuid, "' 错误： ", err.Error(), "\n")
		fmt.Println(text)
		dir.WriteString("error.txt", text)
		panic(err)
	}
	fmt.Println("++++++++ dorp db ++++++++")
	db, err := database.Instance().DB()
	if err != nil {
		zap_server.ZAPLOG.Error("获取数据库连接失败", zap.String("database.Instance().DB()", err.Error()))
		panic(err)
	}
	if db != nil {
		db.Close()
	}

	if multi.AuthDriver != nil {
		multi.AuthDriver.Close()
	}
	err = database.Remove()
	if err != nil {
		zap_server.ZAPLOG.Error("删除配置文件失败", zap.String("database.Remove", err.Error()))
		panic(err)
	}
	if testClient != nil {
		testClient.Logout(logOutUrl, nil)
	}
}

package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/helper/tests"
	"github.com/snowlyg/iris-admin/migration"
	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/operation"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
	"github.com/snowlyg/iris-admin/server/web/web_gin/modules/rbac/admin"
	"github.com/snowlyg/iris-admin/server/web/web_gin/modules/rbac/api"
	"github.com/snowlyg/iris-admin/server/web/web_gin/modules/rbac/authority"
	"github.com/snowlyg/iris-admin/server/web/web_gin/modules/rbac/public"
	multi "github.com/snowlyg/multi/gin"
)

// Party v1 模块
func Party(group *gin.RouterGroup) {
	api.Group(group)
	admin.Group(group)
	authority.Group(group)
	public.Group(group)
}

var LoginResponse = tests.Responses{
	{Key: "status", Value: http.StatusOK},
	{Key: "message", Value: "请求成功"},
	{Key: "data",
		Value: tests.Responses{
			{Key: "accessToken", Value: "", Type: "notempty"},
		},
	},
}
var LogoutResponse = tests.Responses{
	{Key: "status", Value: http.StatusOK},
	{Key: "message", Value: "请求成功"},
}

func BeforeTestMain(mysqlPwd, redisPwd string, redisDB int) (string, *web_gin.WebServer) {
	fmt.Println("+++++ before test +++++")
	node, _ := snowflake.NewNode(1)
	uuid := str.Join("gin", "_", node.Generate().String())
	fmt.Printf("+++++ %s +++++\n\n", uuid)
	web_gin.CONFIG.System.CacheType = "redis"
	web_gin.CONFIG.System.DbType = "mysql"
	web_gin.InitWeb()

	database.CONFIG = database.Mysql{
		Path:         "127.0.0.1:3306",
		Config:       "charset=utf8mb4&parseTime=True&loc=Local",
		Dbname:       uuid,
		Username:     "root",
		Password:     strings.TrimSpace(mysqlPwd),
		MaxIdleConns: 0,
		MaxOpenConns: 0,
		LogMode:      true,
		LogZap:       "error",
	}
	database.InitMysql()

	cache.CONFIG = cache.Redis{
		DB:       redisDB,
		Addr:     "127.0.0.1:6379",
		Password: strings.TrimSpace(redisPwd),
		PoolSize: 0,
	}

	cache.InitCache()

	wi := web_gin.Init()
	Party(wi.GetRouterGroup("/api/v1"))
	web.StartTest(wi)

	mc := migration.New()
	// 添加 v1 内置模块数据表和数据
	fmt.Println("++++++ add model ++++++")
	mc.AddModel(&api.Api{}, &authority.Authority{}, &admin.Admin{}, &operation.Oplog{})
	routes, _ := wi.GetSources()
	fmt.Println("+++++++ seed data ++++++")
	// notice : 注意模块顺序
	mc.AddSeed(api.New(routes), authority.Source, admin.Source)
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

func AfterTestMain(uuid string) {
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
		dir.WriteString("error.txt", err.Error())
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
		dir.WriteString("error.txt", err.Error())
		panic(err)
	}
}

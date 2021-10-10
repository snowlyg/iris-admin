package test

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/modules/migration"
	v1 "github.com/snowlyg/iris-admin/modules/v1"
	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_iris"
	"github.com/snowlyg/multi"
)

//go:embed mysqlPwd.txt
var mysqlPwd string

//go:embed redisPwd.txt
var redisPwd string

var TestServer *web_iris.WebServer

func TestMain(m *testing.M) {
	uuid := uuid.New().String()

	web_iris.CONFIG.System.CacheType = "redis"
	web_iris.CONFIG.System.DbType = "mysql"
	web_iris.InitWeb()

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
		DB:       1,
		Addr:     "127.0.0.1:6379",
		Password: strings.TrimSpace(redisPwd),
		PoolSize: 0,
	}
	cache.InitCache()
	migration.Gormigrate().Migrate()

	TestServer = web_iris.Init()
	v1Party := web_iris.Party{
		Perfix:    "/api/v1",
		PartyFunc: v1.Party(),
	}
	TestServer.AddModule(v1Party)
	web.StartTest(TestServer)

	code := m.Run()

	err := database.DorpDB(database.CONFIG.BaseDsn(), "mysql", uuid)
	if err != nil {
		text := str.Join("删除数据库 '", uuid, "' 错误： ", err.Error(), "\n")
		fmt.Println(text)
		dir.WriteString("error.txt", text)
		panic(err)
	}

	db, _ := database.Instance().DB()
	if db != nil {
		db.Close()
	}

	if multi.AuthDriver != nil {
		multi.AuthDriver.Close()
	}

	os.Exit(code)
}

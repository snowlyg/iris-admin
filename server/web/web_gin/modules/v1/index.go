package v1

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/migration"
	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/operation"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
	"github.com/snowlyg/iris-admin/server/web/web_gin/modules/v1/perm"
	multi "github.com/snowlyg/multi/gin"
)

// Party v1 模块
func Party(group *gin.RouterGroup) {
	perm.Group(group)
	// v1.PartyFunc("/users", user.Party())
	// v1.PartyFunc("/roles", role.Party())
	// v1.PartyFunc("/file", file.Party())
	// v1.PartyFunc("/auth", auth.Party())
	// v1.PartyFunc("/oplog", oplog.Party())
}

func BeforeTestMain(mysqlPwd, redisPwd string, redisDB int) (string, *web_gin.WebServer) {
	uuid := uuid.New().String()

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
	web.StartTest(wi)

	mc := migration.New()
	// 添加 v1 内置模块数据表和数据
	mc.AddModel(&perm.Api{}, &operation.Oplog{})
	routes, _ := wi.GetSources()

	// notice : 注意模块顺序
	mc.AddSeed(perm.New(routes))
	err := mc.Migrate()
	if err != nil {
		panic(err)
	}
	err = mc.Seed()
	if err != nil {
		panic(err)
	}

	return uuid, wi
}

func AfterTestMain(uuid string) {
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

}

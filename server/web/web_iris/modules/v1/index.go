package v1

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/migration"
	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/operation"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_iris"
	"github.com/snowlyg/iris-admin/server/web/web_iris/modules/v1/auth"
	"github.com/snowlyg/iris-admin/server/web/web_iris/modules/v1/file"
	"github.com/snowlyg/iris-admin/server/web/web_iris/modules/v1/oplog"
	"github.com/snowlyg/iris-admin/server/web/web_iris/modules/v1/perm"
	"github.com/snowlyg/iris-admin/server/web/web_iris/modules/v1/role"
	"github.com/snowlyg/iris-admin/server/web/web_iris/modules/v1/user"
	"github.com/snowlyg/multi"
)

// Party v1 模块
func Party() func(v1 iris.Party) {
	return func(v1 iris.Party) {
		v1.PartyFunc("/users", user.Party())
		v1.PartyFunc("/roles", role.Party())
		v1.PartyFunc("/perms", perm.Party())
		v1.PartyFunc("/file", file.Party())
		v1.PartyFunc("/auth", auth.Party())
		v1.PartyFunc("/oplog", oplog.Party())
	}
}

func BeforeTestMain(mysqlPwd, redisPwd string, redisDB int) (string, *web_iris.WebServer) {
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
		DB:       redisDB,
		Addr:     "127.0.0.1:6379",
		Password: strings.TrimSpace(redisPwd),
		PoolSize: 0,
	}

	cache.InitCache()

	wi := web_iris.Init()
	v1Party := web_iris.Party{
		Perfix:    "/api/v1",
		PartyFunc: Party(),
	}
	wi.AddModule(v1Party)
	web.StartTest(wi)

	mc := migration.New()
	// 添加 v1 内置模块数据表和数据
	mc.AddModel(&perm.Permission{}, &role.Role{}, &user.User{}, &operation.Oplog{})
	routes, _ := wi.GetSources()
	// notice : 注意模块顺序
	mc.AddSeed(perm.New(routes), role.Source, user.Source)
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

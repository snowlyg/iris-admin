package tests

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/migration"
	"github.com/snowlyg/iris-admin/server/casbin"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/operation"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
	"github.com/snowlyg/iris-admin/server/web/web_iris"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
)

func BeforeTestMainGin(party func(wi *web_gin.WebServer), seed func(wi *web_gin.WebServer, mc *migration.MigrationCmd)) (string, *web_gin.WebServer) {
	fmt.Println("+++++ before test +++++")

	dbType := os.Getenv("dbType")
	if dbType != "" {
		web.CONFIG.System.DbType = dbType
	}

	web.InitWeb()

	mysqlPwd := os.Getenv("mysqlPwd")
	mysqlAddr := os.Getenv("mysqlAddr")
	if mysqlAddr == "" {
		database.CONFIG.Path = strings.TrimSpace(mysqlAddr)
	}
	node, _ := snowflake.NewNode(1)
	uuid := str.Join("gin", "_", node.Generate().String())

	fmt.Printf("+++++ %s +++++\n\n", uuid)
	database.CONFIG.Dbname = uuid
	database.CONFIG.Password = strings.TrimSpace(mysqlPwd)
	database.CONFIG.LogMode = true
	database.InitMysql()

	wi := web_gin.Init()
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

func BeforeTestMainIris(party func(wi *web_iris.WebServer), seed func(wi *web_iris.WebServer, mc *migration.MigrationCmd)) (string, *web_iris.WebServer) {
	fmt.Println("+++++ before test +++++")

	dbType := os.Getenv("dbType")
	if dbType != "" {
		web.CONFIG.System.DbType = dbType
	}

	web.InitWeb()

	mysqlPwd := os.Getenv("mysqlPwd")
	mysqlAddr := os.Getenv("mysqlAddr")
	if mysqlAddr != "" {
		database.CONFIG.Path = strings.TrimSpace(mysqlAddr)
	}
	node, _ := snowflake.NewNode(1)
	uuid := str.Join("iris", "_", node.Generate().String())

	fmt.Printf("+++++ %s +++++\n\n", uuid)
	database.CONFIG.Dbname = uuid
	database.CONFIG.Password = strings.TrimSpace(mysqlPwd)
	database.CONFIG.LogMode = true
	database.InitMysql()

	wi := web_iris.Init()
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

func AfterTestMain(uuid string, isDelDb bool) {
	fmt.Println("++++++++ after test main ++++++++")
	if isDelDb {
		err := database.DorpDB(database.CONFIG.BaseDsn(), "mysql", uuid)
		if err != nil {
			zap_server.ZAPLOG.Error("删除数据库失败", zap.String("uuid", uuid), zap.String("err", err.Error()))
		}
	}
	fmt.Println("++++++++ dorp db ++++++++")
	db, err := database.Instance().DB()
	if err != nil {
		zap_server.ZAPLOG.Error(str.Join("获取数据库连接失败:", err.Error()))
	}
	if db != nil {
		db.Close()
	}

	defer zap_server.Remove()
	defer operation.Remove()
	defer casbin.Remove()
	defer web.Remove()
	defer database.Remove()
}

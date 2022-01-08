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
	fmt.Println("+++++ 测试开始 +++++")

	dbType := os.Getenv("dbType")
	if dbType != "" {
		web.CONFIG.System.DbType = dbType
	}

	web.InitWeb()

	node, _ := snowflake.NewNode(1)
	uuid := str.Join("gin", "_", node.Generate().String())

	database.CONFIG.Dbname = uuid
	mysqlPwd := os.Getenv("mysqlPwd")
	if mysqlPwd != "" {
		database.CONFIG.Password = strings.TrimSpace(mysqlPwd)
	}
	mysqlAddr := os.Getenv("mysqlAddr")
	if mysqlAddr != "" {
		database.CONFIG.Path = strings.TrimSpace(mysqlAddr)
	}
	database.CONFIG.LogMode = true
	database.InitMysql()

	wi := web_gin.Init()
	party(wi)
	web.StartTest(wi)

	mc := migration.New()
	// 添加 v1 内置模块数据表和数据
	fmt.Println("++++++ 添加数据 ++++++")
	seed(wi, mc)
	err := mc.Migrate()
	if err != nil {
		fmt.Printf("数据库迁移失败： [%s]", err.Error())
		return uuid, nil
	}
	err = mc.Seed()
	if err != nil {
		fmt.Printf("数据填充失败：[%s]", err.Error())
		return uuid, nil
	}

	return uuid, wi
}

func BeforeTestMainIris(party func(wi *web_iris.WebServer), seed func(wi *web_iris.WebServer, mc *migration.MigrationCmd)) (string, *web_iris.WebServer) {
	fmt.Println("+++++ 测试前置方法 +++++")

	dbType := os.Getenv("dbType")
	if dbType != "" {
		web.CONFIG.System.DbType = dbType
	}

	web.InitWeb()

	node, _ := snowflake.NewNode(1)
	uuid := str.Join("iris", "_", node.Generate().String())

	database.CONFIG.Dbname = uuid
	mysqlPwd := os.Getenv("mysqlPwd")
	if mysqlPwd != "" {
		database.CONFIG.Password = strings.TrimSpace(mysqlPwd)
	}
	mysqlAddr := os.Getenv("mysqlAddr")
	if mysqlAddr != "" {
		database.CONFIG.Path = strings.TrimSpace(mysqlAddr)
	}
	database.CONFIG.LogMode = true
	database.InitMysql()

	wi := web_iris.Init()
	party(wi)
	web.StartTest(wi)

	mc := migration.New()

	// 添加 v1 内置模块数据表和数据
	fmt.Println("++++++ 添加数据 ++++++")

	seed(wi, mc)
	err := mc.Migrate()
	if err != nil {
		fmt.Printf("数据库迁移失败： [%s]", err.Error())
		return uuid, nil
	}
	err = mc.Seed()
	if err != nil {
		fmt.Printf("数据填充失败：[%s]", err.Error())
		return uuid, nil
	}

	return uuid, wi
}

func AfterTestMain(uuid string, isDelDb bool) {
	fmt.Println("++++++++ 测试后置方法 ++++++++")
	if isDelDb {
		err := database.DorpDB(database.CONFIG.BaseDsn(), "mysql", uuid)
		if err != nil {
			zap_server.ZAPLOG.Error("删除数据库失败", zap.String("uuid", uuid), zap.String("err", err.Error()))
		}
	}
	fmt.Println("++++++++ 删除数据表 ++++++++")

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

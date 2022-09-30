package common

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

// BeforeTestMainGin
func BeforeTestMainGin(party func(wi *web_gin.WebServer), seed func(wi *web_gin.WebServer, mc *migration.MigrationCmd)) (string, *web_gin.WebServer) {
	fmt.Println("+++++ TEST BEGAIN +++++")

	dbType := os.Getenv("dbType")
	if dbType != "" {
		web.CONFIG.System.DbType = dbType
	}
	web.InitWeb()

	zap_server.CONFIG.LogInConsole = true

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

	// database.InitMysql()

	wi := web_gin.Init()
	party(wi)
	web.StartTest(wi)

	mc := migration.New()
	fmt.Println("++++++ add data to database ++++++")
	seed(wi, mc)
	err := mc.Migrate()
	if err != nil {
		fmt.Printf("migrate fail: [%s]", err.Error())
		return uuid, nil
	}
	err = mc.Seed()
	if err != nil {
		fmt.Printf("seed fail: [%s]", err.Error())
		return uuid, nil
	}

	return uuid, wi
}

// BeforeTestMainIris
func BeforeTestMainIris(party func(wi *web_iris.WebServer), seed func(wi *web_iris.WebServer, mc *migration.MigrationCmd)) (string, *web_iris.WebServer) {
	fmt.Println("+++++ TEST BEGAIN +++++")

	dbType := os.Getenv("dbType")
	if dbType != "" {
		web.CONFIG.System.DbType = dbType
	}
	web.InitWeb()

	zap_server.CONFIG.LogInConsole = true

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

	// database.InitMysql()

	wi := web_iris.Init()
	party(wi)
	web.StartTest(wi)

	mc := migration.New()

	fmt.Println("++++++ add datas to database ++++++")

	seed(wi, mc)
	err := mc.Migrate()
	if err != nil {
		fmt.Printf("migrate fail: [%s]", err.Error())
		return uuid, nil
	}
	err = mc.Seed()
	if err != nil {
		fmt.Printf("seed fail: [%s]", err.Error())
		return uuid, nil
	}

	return uuid, wi
}

func AfterTestMain(uuid string, isDelDb bool) {
	defer fmt.Println("++++++++ AFTER END ++++++++")
	if isDelDb {
		err := database.DorpDB(database.CONFIG.BaseDsn(), "mysql", uuid)
		if err != nil {
			zap_server.ZAPLOG.Error("delete database fail", zap.String("uuid", uuid), zap.String("err", err.Error()))
		}
	}
	fmt.Println("++++++++ DELETE DATABASE ++++++++")

	db, err := database.Instance().DB()
	if err != nil {
		zap_server.ZAPLOG.Error(str.Join("get database instance fail:", err.Error()))
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

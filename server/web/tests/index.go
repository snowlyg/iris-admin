package tests

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/migration"
	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
	"github.com/snowlyg/iris-admin/server/web/web_iris"
	"github.com/snowlyg/iris-admin/server/zap_server"
)

func BeforeTestMainGin(party func(wi *web_gin.WebServer), seed func(wi *web_gin.WebServer, mc *migration.MigrationCmd)) (string, *web_gin.WebServer) {
	fmt.Println("+++++ before test +++++")
	mysqlPwd := os.Getenv("mysqlPwd")
	redisPwd := os.Getenv("redisPwd")
	if strings.TrimSpace(mysqlPwd) != database.CONFIG.Password {
		err := database.Remove()
		if err != nil {
			zap_server.ZAPLOG.Error(str.Join("删除数据库配置文件失败:", err.Error()))
		}
	}
	if strings.TrimSpace(redisPwd) != cache.CONFIG.Password {
		err := cache.Remove()
		if err != nil {
			zap_server.ZAPLOG.Error(str.Join("删除缓存配置文件失败:", err.Error()))
		}
	}
	node, _ := snowflake.NewNode(1)
	uuid := str.Join("gin", "_", node.Generate().String())
	fmt.Printf("+++++ %s +++++\n\n", uuid)
	web.CONFIG.System.DbType = "mysql"
	web.InitWeb()

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
	mysqlPwd := os.Getenv("mysqlPwd")
	redisPwd := os.Getenv("redisPwd")
	if strings.TrimSpace(mysqlPwd) != database.CONFIG.Password {
		err := database.Remove()
		if err != nil {
			zap_server.ZAPLOG.Error(str.Join("删除数据库配置文件失败:", err.Error()))
		}
	}
	if strings.TrimSpace(redisPwd) != cache.CONFIG.Password {
		err := cache.Remove()
		if err != nil {
			zap_server.ZAPLOG.Error(str.Join("删除缓存配置文件失败:", err.Error()))
		}
	}
	node, _ := snowflake.NewNode(1)
	uuid := str.Join("iris", "_", node.Generate().String())
	fmt.Printf("+++++ %s +++++\n\n", uuid)
	web.CONFIG.System.DbType = "mysql"
	web.InitWeb()

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
			text := str.Join("删除数据库 '", uuid, "' 错误： ", err.Error(), "\n")
			zap_server.ZAPLOG.Error(text)
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

	err = database.Remove()
	if err != nil {
		zap_server.ZAPLOG.Error(str.Join("删除数据库配置文件失败:", err.Error()))
	}
	err = web.Remove()
	if err != nil {
		zap_server.ZAPLOG.Error(str.Join("删除缓存配置文件失败:", err.Error()))
	}
	err = web.Remove()
	if err != nil {
		zap_server.ZAPLOG.Error(str.Join("删除web配置文件失败:", err.Error()))
	}
}

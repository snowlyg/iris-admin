package common

import (
	"log"

	"github.com/bwmarrin/snowflake"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/migration"
	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
	"github.com/snowlyg/iris-admin/server/web/web_iris"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"gorm.io/gorm"
)

// BeforeTestMainGin
func BeforeTestMainGin(party func(wi *web_gin.WebServer), seed func(wi *web_gin.WebServer, mc *migration.MigrationCmd)) (string, *web_gin.WebServer, error) {
	zap_server.CONFIG.LogInConsole = true
	if err := zap_server.Recover(); err != nil {
		log.Printf("zap recover fail:%s\n", err.Error())
	}

	dbType := g.TestDbType
	if dbType != "" {
		web.CONFIG.System.DbType = dbType
	}
	if dbType == "redis" {
		if err := cache.Recover(); err != nil {
			log.Printf("cache recover fail:%s\n", err.Error())
		}
	}
	if err := web.Recover(); err != nil {
		log.Printf("web recover fail:%s\n", err.Error())
	}

	node, _ := snowflake.NewNode(1)
	uuid := str.Join("gin", "_", node.Generate().String())

	database.CONFIG.DbName = uuid
	if user := g.TestMysqlName; user != "" {
		database.CONFIG.Username = user
	}
	if pwd := g.TestMysqlPwd; pwd != "" {
		database.CONFIG.Password = pwd
	}
	if addr := g.TestMysqlAddr; addr != "" {
		database.CONFIG.Path = addr
	}
	database.CONFIG.LogMode = true
	if err := database.Recover(); err != nil {
		log.Printf("databse recover fail:%s\n", err.Error())
	}

	if database.Instance() == nil {
		return uuid, nil, gorm.ErrInvalidDB
	}

	wi := web_gin.Init()
	party(wi)
	web.StartTest(wi)

	mc := migration.New()
	seed(wi, mc)
	err := mc.Migrate()
	if err != nil {
		// fmt.Printf("migrate fail: [%s]", err.Error())
		return uuid, nil, err
	}
	err = mc.Seed()
	if err != nil {
		// fmt.Printf("seed fail: [%s]", err.Error())
		return uuid, nil, err
	}

	return uuid, wi, nil
}

// BeforeTestMainIris
func BeforeTestMainIris(party func(wi *web_iris.WebServer), seed func(wi *web_iris.WebServer, mc *migration.MigrationCmd)) (string, *web_iris.WebServer) {
	zap_server.CONFIG.LogInConsole = true
	zap_server.Recover()

	dbType := g.TestDbType
	if dbType != "" {
		web.CONFIG.System.DbType = dbType
	}
	if err := web.Recover(); err != nil {
		log.Printf("web config recover fail %s\n", err.Error())
	}

	node, _ := snowflake.NewNode(1)
	uuid := str.Join("iris", "_", node.Generate().String())

	database.CONFIG.DbName = uuid
	if user := g.TestMysqlName; user != "" {
		database.CONFIG.Username = user
	}
	if pwd := g.TestMysqlPwd; pwd != "" {
		database.CONFIG.Password = pwd
	}
	if addr := g.TestMysqlAddr; addr != "" {
		database.CONFIG.Path = addr
	}
	database.CONFIG.LogMode = true

	if err := database.Recover(); err != nil {
		log.Printf("database config recover fail %s\n", err.Error())
	}

	if database.Instance() == nil {
		log.Println("database instance is nil")
		return uuid, nil
	}

	wi := web_iris.Init()
	party(wi)
	web.StartTest(wi)

	mc := migration.New()

	seed(wi, mc)
	if err := mc.Migrate(); err != nil {
		log.Printf("migrate fail: [%s]\n", err.Error())
		return uuid, nil
	}
	if err := mc.Seed(); err != nil {
		log.Printf("seed fail: [%s]\n", err.Error())
		return uuid, nil
	}

	return uuid, wi
}

func AfterTestMain(uuid string, isDelDb bool) {
	if isDelDb {
		dsn := database.CONFIG.BaseDsn()
		if err := database.DorpDB(dsn, "mysql", uuid); err != nil {
			log.Printf("drop table(%s) on dsn(%s) fail %s\n", uuid, dsn, err.Error())
		}
	}

	if db, _ := database.Instance().DB(); db != nil {
		db.Close()
	}
}

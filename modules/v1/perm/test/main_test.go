package test

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"
	v1 "github.com/snowlyg/iris-admin/modules/v1"
	"github.com/snowlyg/iris-admin/modules/v1/initdb"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/multi"
)

//go:embed mysql_pwd.txt
var mysql_pwd string

//go:embed redis_pwd.txt
var redis_pwd string

var TestServer *web.WebServer

func TestMain(m *testing.M) {
	TestServer = web.Init()
	uuid := uuid.NewV3(uuid.NewV4(), uuid.NamespaceOID.String()).String()
	config := initdb.Request{
		SqlType: "mysql",
		Sql: initdb.Sql{
			Host:     "127.0.0.1",
			Port:     "3306",
			UserName: "root",
			Password: strings.TrimSpace(mysql_pwd),
			DBName:   uuid,
			LogMode:  true,
		},
		CacheType: "redis",
		Cache: initdb.Cache{
			Host:     "127.0.0.1",
			Port:     "6379",
			Password: strings.TrimSpace(redis_pwd),
			PoolSize: 1000,
			DB:       2,
		},
		Addr:  "127.0.0.1:8085",
		Level: "test",
	}

	TestServer.AddModule(v1.Party())
	err := TestServer.InitRouter()
	if err != nil {
		fmt.Printf("初始化路由错误： %v\n", err)
		panic(err)
	}
	err = initdb.InitDB(config)
	if err != nil {
		fmt.Printf("初始化数据库错误： %v\n", err)
		panic(err)
	}

	code := m.Run()

	err = dorpDB(uuid)
	if err != nil {
		fmt.Printf("删除数据库 '%s' 错误： %v\n", uuid, err)
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

func dorpDB(uuid string) error {
	if err := database.Instance().Exec(fmt.Sprintf("drop database if exists `%s`;", uuid)).Error; err != nil {
		return err
	}

	return nil
}

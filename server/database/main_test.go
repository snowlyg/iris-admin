package database

import (
	_ "embed"
	"os"
	"strings"
	"testing"

	"github.com/bwmarrin/snowflake"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	mysqlPwd := os.Getenv("mysqlPwd")
	node, _ := snowflake.NewNode(1)
	uuid := str.Join("database", "_", node.Generate().String())

	CONFIG = Mysql{
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
	InitMysql()

	Instance()

	code := m.Run()

	err := DorpDB(CONFIG.BaseDsn(), "mysql", uuid)
	if err != nil {
		zap_server.ZAPLOG.Error("删除数据库失败", zap.String("uuid", uuid), zap.String("err", err.Error()))
		panic(err)
	}

	db, err := Instance().DB()
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		panic(err)
	}
	if db != nil {
		db.Close()
	}

	err = Remove()
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		panic(err)
	}
	zap_server.Remove()
	os.Exit(code)
}

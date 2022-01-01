package casbin

import (
	_ "embed"
	"os"
	"strings"
	"testing"

	"github.com/bwmarrin/snowflake"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/zap_server"
)

func TestMain(m *testing.M) {
	mysqlPwd := os.Getenv("mysqlPwd")
	node, _ := snowflake.NewNode(1)
	uuid := str.Join("casbin", "_", node.Generate().String())

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

	Instance()

	code := m.Run()

	err := database.DorpDB(database.CONFIG.BaseDsn(), "mysql", uuid)
	if err != nil {
		zap_server.ZAPLOG.Error(str.Join("删除数据库 '", uuid, "' 错误： ", err.Error(), "\n"))
		panic(err)
	}

	db, err := database.Instance().DB()
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		panic(err)
	}
	if db != nil {
		db.Close()
	}
	err = database.Remove()
	if err != nil {
		zap_server.ZAPLOG.Error(err.Error())
		panic(err)
	}
	Remove()
	zap_server.Remove()
	os.Exit(code)
}

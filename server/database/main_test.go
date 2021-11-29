package database

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bwmarrin/snowflake"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
)

//go:embed mysqlPwd.txt
var mysqlPwd string

func TestMain(m *testing.M) {
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
		text := str.Join("删除数据库 '", uuid, "' 错误： ", err.Error(), "\n")
		fmt.Println(text)
		dir.WriteString("error.txt", text)
		panic(err)
	}

	db, err := Instance().DB()
	if err != nil {
		dir.WriteString("error.txt", err.Error())
		panic(err)
	}
	if db != nil {
		db.Close()
	}

	err = Remove()
	if err != nil {
		dir.WriteString("error.txt", err.Error())
		panic(err)
	}

	os.Exit(code)
}

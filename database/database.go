package database

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/snowlyg/IrisAdminApi/config"
	"github.com/snowlyg/IrisAdminApi/files"
)

var db *dataBase
var once sync.Once

type dataBase struct {
	Db       *gorm.DB
	Enforcer *casbin.Enforcer
}

/**
*设置数据库连接
*@param diver string
 */
func getDataBase() *dataBase {
	once.Do(func() {
		driverName, conn := getDriverNameAndConn()
		gdb, err := gorm.Open(driverName, conn)
		if err != nil {
			color.Red(fmt.Sprintf("gorm open 错误: %v", err))
		}

		c, err := gormadapter.NewAdapter(driverName, conn, true) // Your driver and data source.
		if err != nil {
			color.Red(fmt.Sprintf("NewAdapter 错误: %v", err))
		}

		e, err := casbin.NewEnforcer(files.GetAbsPath("database", "rbac_model.conf"), c)
		if err != nil {
			color.Red(fmt.Sprintf("NewEnforcer 错误: %v", err))
		}

		_ = e.LoadPolicy()

		db = &dataBase{Db: gdb, Enforcer: e}
	})
	return db
}

/*
	获取数据连接驱动类型和链接
*/
func getDriverNameAndConn() (string, string) {
	var driverName string
	var conn string
	driverType := config.GetAppDriverType()
	if driverType == "Sqlite" {
		driverName = "sqlite3"
		if isTestEnv() {
			conn = config.GetSqliteTConnect()
		} else {
			conn = config.GetSqliteConnect()
		}
	} else if driverType == "Mysql" {
		driverName = "mysql"
		connect := config.GetMysqlConnect()
		if isTestEnv() {
			conn = connect + config.GetMysqlTName() + "?charset=utf8&parseTime=True&loc=Local"
		} else {
			conn = connect + config.GetMysqlName() + "?charset=utf8&parseTime=True&loc=Local"
		}
	}

	return driverName, conn
}

func GetGdb() *gorm.DB {
	return getDataBase().Db
}

func GetEnforcer() *casbin.Enforcer {
	return getDataBase().Enforcer
}

func isTestEnv() bool {
	files := os.Args
	for _, v := range files {
		if strings.Contains(v, "test") {
			return true
		}
	}
	return false
}

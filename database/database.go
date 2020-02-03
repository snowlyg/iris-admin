package database

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"IrisAdminApi/config"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var db *dataBase
var once sync.Once
var path string

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
		dirverName, conn := getDirverNameAndConn()
		gdb, err := gorm.Open(dirverName, conn)
		if err != nil {
			color.Red(fmt.Sprintf("gorm open 错误: %v", err))
		}

		c, err := gormadapter.NewAdapter(dirverName, conn, true) // Your driver and data source.
		if err != nil {
			color.Red(fmt.Sprintf("NewAdapter 错误: %v", err))
		}

		if len(path) == 0 {
			path = "database/rbac_model.conf"
		}

		e, err := casbin.NewEnforcer(path, c)
		if err != nil {
			color.Red(fmt.Sprintf("NewEnforcer 错误: %v", err))
		}

		_ = e.LoadPolicy()

		db = &dataBase{Db: gdb, Enforcer: e}
	})
	return db
}

/*
	重置配置文件路径
	该方法目前仅用于测试 ，有些多余
*/
func SetDatabasePath(dbPath string)  {
	path = dbPath
}

/*
	获取数据连接驱动类型和链接
*/
func getDirverNameAndConn() (string, string) {
	var dirverName string
	var conn string
	dirverType := config.GetAppDirverType()
	if dirverType == "Sqlite" {
		dirverName = "sqlite3"
		if isTestEnv() {
			conn = config.GetSqliteTConnect()
		} else {
			conn = config.GetSqliteConnect()
		}
	} else if dirverType == "Mysql" {
		dirverName = "mysql"
		connect := config.GetMysqlConnect()
		if isTestEnv() {
			conn = connect + config.GetMysqlTName() + "?charset=utf8&parseTime=True&loc=Local"
		} else {
			conn = connect + config.GetMysqlName() + "?charset=utf8&parseTime=True&loc=Local"
		}
	}

	return dirverName, conn
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


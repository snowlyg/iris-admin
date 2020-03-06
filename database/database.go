package database

import (
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/snowlyg/IrisAdminApi/config"
	"github.com/snowlyg/IrisAdminApi/files"
	"github.com/snowlyg/IrisAdminApi/libs"
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
		conn := getDriverConn()
		driverType := config.GetAppDriverType()

		gdb, err := gorm.Open(driverType, conn)
		if err != nil {
			color.Red(fmt.Sprintf("gorm open 错误: %v", err))
		}

		//c, err := gormadapter.NewAdapterByDB(gdb)
		c, err := gormadapter.NewAdapter(driverType, conn, true) // Your driver and data source.
		if err != nil {
			color.Red(fmt.Sprintf("NewAdapter 错误: %v", err))
		}

		e, err := casbin.NewEnforcer(files.GetAbsPath("database", "rbac_model.conf"), c)
		if err != nil {
			color.Red(fmt.Sprintf("NewEnforcer 错误: %v", err))
		}

		// 修改默认匹配
		//rm := defaultrolemanager.NewRoleManager(10).(*defaultrolemanager.RoleManager)
		//rm.AddMatchingFunc("KeyMatch3", util.KeyMatch3)

		_ = e.LoadPolicy()

		db = &dataBase{Db: gdb, Enforcer: e}
	})
	return db
}

//获取数据连接驱动类型和链接
func getDriverConn() string {
	var conn string
	driverType := config.GetAppDriverType()
	if driverType == "sqlite3" {
		if libs.IsTestEnv() {
			conn = config.GetSqliteTConnect()
		} else {
			conn = config.GetSqliteConnect()
		}
	} else if driverType == "mysql" {
		connect := config.GetMysqlConnect()
		if libs.IsTestEnv() {
			conn = connect + config.GetMysqlTName() + "?charset=utf8&parseTime=True&loc=Local"
		} else {
			conn = connect + config.GetMysqlName() + "?charset=utf8&parseTime=True&loc=Local"
		}
	}

	return conn
}

func GetGdb() *gorm.DB {
	return getDataBase().Db
}

func GetEnforcer() *casbin.Enforcer {
	return getDataBase().Enforcer
}

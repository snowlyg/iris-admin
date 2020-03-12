package database

import (
	"fmt"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/snowlyg/IrisAdminApi/config"
	"github.com/snowlyg/IrisAdminApi/libs"
)

var (
	Db       *gorm.DB
	Enforcer *casbin.Enforcer
)

/**
*设置数据库连接
*@param diver string
 */
func init() {
	var err error
	conn := getDriverConn()
	driverType := config.GetAppDriverType()

	Db, err = gorm.Open(driverType, conn)
	if err != nil {
		color.Red(fmt.Sprintf("gorm open 错误: %v", err))
	}

	//c, err := gormadapter.NewAdapterByDB(gdb)
	c, err := gormadapter.NewAdapter(driverType, conn, true) // Your driver and data source.
	if err != nil {
		color.Red(fmt.Sprintf("NewAdapter 错误: %v", err))
	}

	Enforcer, err = casbin.NewEnforcer(filepath.Join(config.Root, "config", "rbac_model.conf"), c)
	if err != nil {
		color.Red(fmt.Sprintf("NewEnforcer 错误: %v", err))
	}

	_ = Enforcer.LoadPolicy()

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

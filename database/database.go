package database

import (
	"fmt"
	"os"
	"strconv"
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

var db *DataBase
var once sync.Once

type DataBase struct {
	Db       *gorm.DB
	Enforcer *casbin.Enforcer
}

/**
*设置数据库连接
*@param diver string
 */
func getDataBase() *DataBase {
	once.Do(func() {
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
				conn = connect +config.GetMysqlName() + "?charset=utf8&parseTime=True&loc=Local"
			}
		}

		gdb, err := gorm.Open(dirverName, conn)
		if err != nil {
			color.Red(fmt.Sprintf("gorm open 错误: %v", err))
		}

		c, err := gormadapter.NewAdapter(dirverName, conn, true) // Your driver and data source.
		if err != nil {
			color.Red(fmt.Sprintf("NewAdapter 错误: %v", err))
		}

		e, err := casbin.NewEnforcer("./config/rbac_model.conf", c)
		if err != nil {
			color.Red(fmt.Sprintf("NewEnforcer 错误: %v", err))
		}
		_ = e.LoadPolicy()
		db = &DataBase{Db: gdb, Enforcer: e}
	})

	return db
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
func Update(v, d interface{}) error {
	if err := GetGdb().Model(v).Updates(d).Error; err != nil {
		return err
	}
	return nil
}

func GetRolesForUser(uid uint) []string {
	uids, err := GetEnforcer().GetRolesForUser(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		color.Red(fmt.Sprintf("GetRolesForUser 错误: %v", err))
		return []string{}
	}

	return uids
}

func GetPermissionsForUser(uid uint) [][]string {
	return GetEnforcer().GetPermissionsForUser(strconv.FormatUint(uint64(uid), 10))
}

func DropTables() {
	GetGdb().DropTable("users", "roles", "permissions", "oauth_tokens", "casbin_rule")
}

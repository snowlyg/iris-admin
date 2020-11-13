package database

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/snowlyg/blog/libs"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"path/filepath"
	"sync"
	"time"

	_ "gorm.io/driver/postgres"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Db       *gorm.DB
	Enforcer *casbin.Enforcer
}

var DB *Database
var loadDatabaseOnce sync.Once

func Singleton() *Database {
	loadDatabaseOnce.Do(getDatabase)
	return DB
}

func getDatabase() {
	DB = new(Database)
	DB.Db = getGormDb()
	DB.Enforcer = getEnforcer()
}

func getConn() string {
	if libs.Config.DB.Adapter == "mysql" {
		return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", libs.Config.DB.User, libs.Config.DB.Password, libs.Config.DB.Host, libs.Config.DB.Port, libs.Config.DB.Name)
	} else if libs.Config.DB.Adapter == "postgres" {
		return fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", libs.Config.DB.User, libs.Config.DB.Password, libs.Config.DB.Host, libs.Config.DB.Name)
	} else if libs.Config.DB.Adapter == "sqlite3" {
		return libs.DBFile()
	} else {
		logger.Println(errors.New("not supported database adapter"))
	}
	return ""
}

// getGormDb
func getGormDb() (Db *gorm.DB) {
	var err error
	var dialector gorm.Dialector
	if libs.Config.DB.Adapter == "mysql" {
		dialector = mysql.Open(getConn())
	} else if libs.Config.DB.Adapter == "postgres" {
		dialector = postgres.Open(getConn())
	} else if libs.Config.DB.Adapter == "sqlite3" {
		dialector = sqlite.Open(getConn())
	} else {
		logger.Println(errors.New("not supported database adapter"))
	}

	Db, err = gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   libs.Config.DB.Prefix, // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: false,                 // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	})

	if err != nil {
		logger.Println(err)
	}

	_ = Db.Use(
		dbresolver.Register(
			dbresolver.Config{ /* xxx */ }).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	)
	Db.Session(&gorm.Session{FullSaveAssociations: true, AllowGlobalUpdate: false})
	return Db
}

// getEnforcer
func getEnforcer() (Enforcer *casbin.Enforcer) {
	var err error

	c, err := gormadapter.NewAdapter(libs.Config.DB.Adapter, getConn(), true) // Your driver and data source.
	if err != nil {
		logger.Println(fmt.Sprintf("NewAdapter 错误: %v,Path: %s", err, getConn()))
	}

	casbinModelPath := filepath.Join(libs.CWD(), "rbac_model.conf")
	Enforcer, err = casbin.NewEnforcer(casbinModelPath, c)
	if err != nil {
		logger.Println(fmt.Sprintf("NewEnforcer 错误: %v", err))
	}

	_ = Enforcer.LoadPolicy()

	return Enforcer
}

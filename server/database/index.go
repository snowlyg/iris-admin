package database

import (
	"os"
	"sync"

	con "github.com/snowlyg/iris-admin/server/config"
	myzap "github.com/snowlyg/iris-admin/server/zap"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	once sync.Once
	db   *gorm.DB
)

// Instance 数据库单例
func Instance() *gorm.DB {
	once.Do(func() {
		switch con.CONFIG.System.DbType {
		case "mysql":
			db = GormMysql()
		default:
			db = GormMysql()
		}
	})
	return db
}

// MysqlTables 注册数据库表专用
func MysqlTables(db *gorm.DB) {
	err := db.AutoMigrate()
	if err != nil {
		myzap.ZAPLOG.Error("注册数据表错误", zap.Any("err", err))
		os.Exit(0)
	}
	myzap.ZAPLOG.Info("注册数据表成功")
}

// GormMysql 初始化Mysql数据库
func GormMysql() *gorm.DB {
	m := con.CONFIG.Mysql
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		DisableDatetimePrecision:  true,    // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,    // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(m.LogMode)); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}

// gormConfig 根据配置决定是否开启日志
func gormConfig(mod bool) *gorm.Config {
	var config = &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	switch con.CONFIG.Mysql.LogZap {
	case "silent", "Silent":
		config.Logger = Default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = Default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = Default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = Default.LogMode(logger.Info)
	case "zap", "Zap":
		config.Logger = Default.LogMode(logger.Info)
	default:
		if mod {
			config.Logger = Default.LogMode(logger.Info)
			break
		}
		config.Logger = Default.LogMode(logger.Silent)
	}
	return config
}

// InitDBFunc 数据化初始化接口
type InitDBFunc interface {
	Init() (err error)
}

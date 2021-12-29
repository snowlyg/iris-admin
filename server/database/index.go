package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/snowlyg/iris-admin/server/viper_server"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ErrDatabaseInit = errors.New("数据库初始化失败")

var (
	once sync.Once
	db   *gorm.DB
)

// InitMysql 初始化数据库
func InitMysql() {
	viper_server.Init(getViperConfig())
}

// Instance 数据库单例
func Instance() *gorm.DB {
	once.Do(func() {
		InitMysql()
		db = gormMysql()
	})
	return db
}

// gormMysql 初始化Mysql数据库
func gormMysql() *gorm.DB {
	if CONFIG.Dbname == "" {
		fmt.Println("config dbname is empty")
		return nil
	}
	err := createTable(CONFIG.BaseDsn(), "mysql", CONFIG.Dbname)
	if err != nil {
		fmt.Printf("create database %s is failed %v \n", CONFIG.Dbname, err)
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       CONFIG.Dsn(), // DSN data source name
		DefaultStringSize:         191,          // string 类型字段的默认长度
		DisableDatetimePrecision:  true,         // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,         // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,         // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,        // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(CONFIG.LogMode)); err != nil {
		fmt.Printf("open mysql is failed %v \n", err)
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(CONFIG.MaxIdleConns)
		sqlDB.SetMaxOpenConns(CONFIG.MaxOpenConns)
		return db
	}
}

// gormConfig 根据配置决定是否开启日志
func gormConfig(mod bool) *gorm.Config {
	var config = &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	switch CONFIG.LogZap {
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

// createTable 创建数据库(mysql)
func createTable(dsn, driver, dbName string) error {
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", dbName)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
}

func DorpDB(dsn, driver, dbName string) error {
	execSql := fmt.Sprintf("DROP database if exists `%s`;", dbName)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	if db == nil {
		return errors.New("db is nil")
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(execSql)
	if err != nil {
		return err
	}
	fmt.Println(execSql)
	return nil
}

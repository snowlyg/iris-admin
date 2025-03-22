package admin

import (
	"database/sql"
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ErrDatabaseInit = errors.New("database initialize fail")

var (
	// once sync.Once
	db *gorm.DB
)

// // init
// func init() {
// 	admin.Init(getViperConfig())
// }

// // Instance
// func Instance() *gorm.DB {
// 	once.Do(func() {
// 		db = gormMysql()
// 	})
// 	return db
// }

// gormMysql get *gorm.DB
func (ws *WebServe) gormMysql() *gorm.DB {
	if ws.conf.Mysql.DbName == "" {
		fmt.Println("config dbname is empty")
		return nil
	}
	err := createTable(ws.conf.Mysql.BaseDsn(), "mysql", ws.conf.Mysql.DbName)
	if err != nil {
		fmt.Printf("create database %s is failed %v \n", ws.conf.Mysql.DbName, err)
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       ws.conf.Mysql.Dsn(),
		DefaultStringSize:         191,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(ws.conf.Mysql.LogMode)); err != nil {
		fmt.Printf("open mysql is failed %v \n", err)
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(ws.conf.Mysql.MaxIdleConns)
		sqlDB.SetMaxOpenConns(ws.conf.Mysql.MaxOpenConns)
		return db
	}
}

// gormConfig get gorm config
func gormConfig(mod bool) *gorm.Config {
	var config = &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	// switch CONFIG.LogZap {
	// case "silent", "Silent":
	// 	config.Logger = Default.LogMode(logger.Silent)
	// case "error", "Error":
	// 	config.Logger = Default.LogMode(logger.Error)
	// case "warn", "Warn":
	// 	config.Logger = Default.LogMode(logger.Warn)
	// case "info", "Info":
	// 	config.Logger = Default.LogMode(logger.Info)
	// case "zap", "Zap":
	// 	config.Logger = Default.LogMode(logger.Info)
	// default:
	// 	if mod {
	// 		config.Logger = Default.LogMode(logger.Info)
	// 		break
	// 	}
	// 	config.Logger = Default.LogMode(logger.Silent)
	// }
	return config
}

// createTable create database(mysql)
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

func dorpDB(dsn, driver, dbName string) error {
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
	return nil
}

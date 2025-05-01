package admin

import (
	"database/sql"
	"fmt"

	"github.com/snowlyg/iris-admin/conf"
	"github.com/snowlyg/iris-admin/e"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// gormDb
func gormDb(conf *conf.Mysql) (*gorm.DB, error) {
	if conf == nil {
		return nil, e.ErrConfigInvalid
	}
	if conf.DbName == "" {
		return nil, e.ErrDbTableNameEmpty
	}
	if err := createTable(conf.BaseDsn(), "mysql", conf.DbName); err != nil {
		return nil, fmt.Errorf("create database %s is fail:%w", conf.DbName, err)
	}
	mysqlConfig := mysql.Config{
		DSN:               conf.Dsn(),
		DefaultStringSize: 191,
		// DisableDatetimePrecision:  true,
		// DontSupportRenameIndex:    true,
		// DontSupportRenameColumn:   true,
		// SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig)); err != nil {
		fmt.Printf("open mysql[%s] is fail:%v\n", conf.Dsn(), err)
		return nil, err
	} else {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
		return db, nil
	}
}

// createTable
func createTable(dsn, driver, dbName string) error {
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", dbName)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("open mysql[%s] is fail:%w", dsn, err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	if err = db.Ping(); err != nil {
		return fmt.Errorf("ping mysql[%s] is fail:%w", dsn, err)
	}
	_, err = db.Exec(createSql)
	return err
}

// func dorpDB(dsn, driver, dbName string) error {
// 	execSql := fmt.Sprintf("DROP database if exists `%s`;", dbName)
// 	db, err := sql.Open(driver, dsn)
// 	if err != nil {
// 		return err
// 	}
// 	if db == nil {
// 		return errors.New("db is nil")
// 	}
// 	defer func(db *sql.DB) {
// 		_ = db.Close()
// 	}(db)
// 	if err = db.Ping(); err != nil {
// 		return err
// 	}
// 	_, err = db.Exec(execSql)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

package admin

import (
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
	// if err := createTable(ws.Config().BaseDsn(), "mysql",conf.DbName); err != nil {
	// 	fmt.Printf("create database %s is failed %v \n",conf.DbName, err)
	// 	return nil
	// }
	mysqlConfig := mysql.Config{
		DSN:                       conf.Dsn(),
		DefaultStringSize:         191,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig)); err != nil {
		fmt.Printf("open mysql is failed %v \n", err)
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

// // createTable create database(mysql)
// func createTable(dsn, driver, dbName string) error {
// 	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", dbName)
// 	db, err := sql.Open(driver, dsn)
// 	if err != nil {
// 		return err
// 	}
// 	defer func(db *sql.DB) {
// 		_ = db.Close()
// 	}(db)
// 	if err = db.Ping(); err != nil {
// 		return err
// 	}
// 	_, err = db.Exec(createSql)
// 	return err
// }

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

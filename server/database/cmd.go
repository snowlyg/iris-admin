package database

import (
	"fmt"
	"strings"
)

// Init 初始化 mysql 配置
func Init() error {
	var cover string
	if IsExist() {
		fmt.Println("Your database config is initialized , reinitialized database will cover your database config.")
		fmt.Println("Did you want to do it ?  [Y/N]")
		fmt.Scanln(&cover)
		switch strings.ToUpper(cover) {
		case "Y":
			err := Remove()
			if err != nil {
				return err
			}
		case "N":
			return nil
		default:
		}
	}
	err := initConfig()
	if err != nil {
		return err
	}
	fmt.Println("mysql initialized finished!")
	return nil
}

func initConfig() error {
	var dbPath, dbName, dbUsername, dbPwd, dbLogZap, dbLogMod string
	var maxIdleConns, maxOpenConns int
	fmt.Println("Please input your database path: ")
	fmt.Println("Database path default is 127.0.0.1:3306")
	fmt.Scanln(&dbPath)

	if dbPath == "" {
		dbPath = "127.0.0.1:3306"
	}
	CONFIG.Path = dbPath

	fmt.Println("Please input your database db-name: ")
	fmt.Println("Database db-name default is iris-admin")
	fmt.Scanln(&dbName)
	if dbName == "" {
		dbName = "iris-admin"
	}
	CONFIG.Dbname = dbName

	fmt.Println("Please input your database username: ")
	fmt.Println("Database username default is root")
	fmt.Scanln(&dbUsername)
	if dbUsername == "" {
		dbUsername = "root"
	}
	CONFIG.Username = dbUsername

	fmt.Println("Please input your database password: ")
	fmt.Scanln(&dbPwd)
	if dbPwd == "" {
		dbPwd = ""
	}
	CONFIG.Password = dbPwd

	fmt.Println("Please input your database log zap: ")
	fmt.Println("Database log zap default is error")
	fmt.Scanln(&dbLogZap)
	if dbLogZap == "" {
		dbLogZap = "error"
	}
	CONFIG.LogZap = dbLogZap

	fmt.Println("Please input your database log mode: [Y/N]")
	fmt.Println("Database log mode default is N")
	fmt.Scanln(&dbLogMod)
	switch strings.ToUpper(dbLogMod) {
	case "Y":
		CONFIG.LogMode = true
	case "N":
		CONFIG.LogMode = false
	default:
		CONFIG.LogMode = false
	}

	fmt.Println("Please input your database max idle conns: ")
	fmt.Scanln(&maxIdleConns)
	CONFIG.MaxIdleConns = maxIdleConns

	fmt.Println("Please input your database max open conns: ")
	fmt.Scanln(&maxOpenConns)
	CONFIG.MaxOpenConns = maxOpenConns

	if Instance() == nil {
		return ErrDatabaseNotInit
	}
	return nil
}

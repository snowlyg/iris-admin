package database

import (
	"fmt"
	"strings"
)

// Init initialize mysql config file
func Init() error {
	var cover string
	if IsExist() {
		fmt.Println("Your database config is initialized , reinitialized database will cover your database config.")
		fmt.Println("Did you want to do it ?  [Y/N]")
		fmt.Scanln(&cover)
		switch strings.ToUpper(cover) {
		case "Y":
		case "N":
			return nil
		default:
		}
	}
	err := Remove()
	if err != nil {
		return err
	}

	err = initConfig()
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
	fmt.Printf("Database path default is '%s'\n", CONFIG.Path)
	fmt.Scanln(&dbPath)
	if dbPath != "" {
		CONFIG.Path = dbPath
	}

	fmt.Println("Please input your database db-name: ")
	fmt.Printf("Database db-name default is '%s'\n", CONFIG.DbName)
	fmt.Scanln(&dbName)
	if dbName != "" {
		CONFIG.DbName = dbName
	}

	fmt.Println("Please input your database username: ")
	fmt.Printf("Database username default is '%s'\n", CONFIG.Username)
	fmt.Scanln(&dbUsername)
	if dbUsername != "" {
		CONFIG.Username = dbUsername
	}

	fmt.Println("Please input your database password: ")
	fmt.Printf("Database password default is '%s'\n", CONFIG.Password)
	fmt.Scanln(&dbPwd)
	if dbPwd != "" {
		CONFIG.Password = dbPwd
	}

	fmt.Println("Please input your database log zap: ")
	fmt.Printf("Database log zap default is '%s'\n", CONFIG.LogZap)
	fmt.Scanln(&dbLogZap)
	if dbLogZap != "" {
		CONFIG.LogZap = dbLogZap
	}

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
	if maxIdleConns > 0 {
		CONFIG.MaxIdleConns = maxIdleConns
	}

	fmt.Println("Please input your database max open conns: ")
	fmt.Scanln(&maxOpenConns)
	if maxOpenConns > 0 {
		CONFIG.MaxOpenConns = maxOpenConns
	}

	if Instance() == nil {
		return ErrDatabaseInit
	}

	return nil
}

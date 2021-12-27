package web

import (
	"fmt"
	"strings"

	"github.com/snowlyg/iris-admin/server/database"
)

// Init 初始化配置文件
func Init() error {
	var cover string
	if IsExist() {
		fmt.Println("Your web config is initialized , reinitialized web will cover your web config.")
		fmt.Println("Did you want to do it ?  [Y/N]")
		fmt.Scanln(&cover)
		switch strings.ToUpper(cover) {
		case "Y":
			err := Remove()
			if err != nil {
				return err
			}
			return initConfig()
		case "N":
			return nil
		default:
			return nil
		}
	}
	err := initConfig()
	if err != nil {
		return err
	}
	fmt.Println("web iris initialized finished!")
	return nil
}

func initConfig() error {
	var dbType string
	fmt.Println("Please choose your database type: ")
	fmt.Println("1. mysql (only support mysql now)")
	fmt.Scanln(&dbType)
	switch dbType {
	case "1":
		CONFIG.System.DbType = "mysql"
		if err := database.Init(); err != nil {
			return err
		}
	default:
		CONFIG.System.DbType = "mysql"
		if err := database.Init(); err != nil {
			return err
		}
	}

	var systemStaticPrefix, systemWebPrefix, systemTimeFormat, systemStaticAbsPath, systemAddr string
	fmt.Println("Please input your system static prefix: ")
	fmt.Println("System static prefix is ''")
	fmt.Scanln(&systemStaticPrefix)
	CONFIG.System.StaticPrefix = systemStaticPrefix

	fmt.Println("Please input your system web prefix: ")
	fmt.Println("System web prefix is ''")
	fmt.Scanln(&systemWebPrefix)
	CONFIG.System.WebPrefix = systemWebPrefix

	fmt.Println("Please input your system timeformat: ")
	fmt.Println("System timeformat is '2006-01-02 15:04:05'")
	fmt.Scanln(&systemTimeFormat)
	if systemTimeFormat == "" {
		systemTimeFormat = "2006-01-02 15:04:05"
	}
	CONFIG.System.TimeFormat = systemTimeFormat

	fmt.Println("Please input your system static abs path: ")
	fmt.Println("System static abs path is ''")
	fmt.Scanln(&systemStaticAbsPath)
	CONFIG.System.StaticAbsPath = systemStaticAbsPath

	fmt.Println("Please input your system addr: ")
	fmt.Println("System addr is '127.0.0.1:8085'")
	fmt.Scanln(&systemAddr)
	if systemAddr == "" {
		systemAddr = "127.0.0.1:8085"
	}
	CONFIG.System.Addr = systemAddr
	err := InitWeb()
	if err != nil {
		return err
	}
	return nil
}

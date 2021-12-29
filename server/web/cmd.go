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
	fmt.Println("web iris-admin initialized finished!")
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

	var systemTimeFormat, systemAddr string
	fmt.Println("Please input your system timeformat: ")
	fmt.Printf("System timeformat is '%s'\n", CONFIG.System.TimeFormat)
	fmt.Scanln(&systemTimeFormat)
	if systemTimeFormat != "" {
		CONFIG.System.TimeFormat = systemTimeFormat
	}

	fmt.Println("Please input your system addr: ")
	fmt.Printf("System addr is '%s'\n", CONFIG.System.Addr)
	fmt.Scanln(&systemAddr)
	if systemAddr != "" {
		CONFIG.System.Addr = systemAddr
	}
	err := InitWeb()
	if err != nil {
		return err
	}
	return nil
}

package cache

import (
	"fmt"
	"strings"
)

// InitConfig 初始化 redis 配置
func InitConfig() error {
	var cover string
	if IsExist() {
		fmt.Println("Your redis config is initialized , reinitialized redis will cover your redis config.")
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

	fmt.Println("redis initialized finished!")
	return nil
}

func initConfig() error {
	var addr, dbPwd string
	var db, poolSize int
	fmt.Println("Please input your redis addr: ")
	fmt.Println("Redis addr default is 127.0.0.1:6379")
	fmt.Scanln(&addr)

	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	CONFIG.Addr = addr

	fmt.Println("Please input your redis db: ")
	fmt.Println("Redis db default is 0")
	fmt.Scanln(&db)
	CONFIG.DB = db

	fmt.Println("Please input your redis password: ")
	fmt.Scanln(&dbPwd)
	if dbPwd == "" {
		dbPwd = ""
	}
	CONFIG.Password = dbPwd

	fmt.Println("Please input your redis pool size: ")
	fmt.Scanln(&poolSize)
	CONFIG.PoolSize = poolSize
	return nil
}

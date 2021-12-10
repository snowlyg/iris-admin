package web_iris

import (
	"fmt"
	"strings"

	"github.com/snowlyg/iris-admin/server/cache"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/multi"
	multi_iris "github.com/snowlyg/multi/iris"
)

// InitConfig 初始化配置文件
func InitConfig() error {
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

			var dbType, cacheType string
			fmt.Println("Please choose your database type: ")
			fmt.Println("1. mysql (only support mysql now)")
			fmt.Scanln(&dbType)
			switch dbType {
			case "1":
				CONFIG.System.DbType = "mysql"
				if err := database.InitConfig(); err != nil {
					return err
				}
			default:
				CONFIG.System.DbType = "mysql"
				if err := database.InitConfig(); err != nil {
					return err
				}
			}

			fmt.Println("Please choose your cache type: ")
			fmt.Println("1. local")
			fmt.Println("2. redis")
			fmt.Scanln(&cacheType)
			switch cacheType {
			case "1":
				CONFIG.System.CacheType = "local"
			case "2":
				CONFIG.System.CacheType = "redis"
				if err := cache.InitConfig(); err != nil {
					return err
				}
			default:
				CONFIG.System.CacheType = "local"
			}

			err = multi_iris.InitDriver(
				&multi.Config{
					DriverType:      CONFIG.System.CacheType,
					UniversalClient: cache.Instance()},
			)
			if err != nil {
				return fmt.Errorf("initialize auth driver failed %w", err)
			}
			if multi.AuthDriver == nil {
				return ErrAuthDriverEmpty
			}
			InitWeb()
		case "N":
			return nil
		default:
		}
	}
	fmt.Println("web iris initialized finished!")
	return nil
}

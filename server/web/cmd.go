package web

import (
	"fmt"
	"strings"

	"github.com/snowlyg/iris-admin/server/database"
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

			var dbType string
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

			InitWeb()
		case "N":
			return nil
		default:
		}
	}
	fmt.Println("web iris initialized finished!")
	return nil
}

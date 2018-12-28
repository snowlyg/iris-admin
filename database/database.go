package database

import (
	"IrisApiProject/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pelletier/go-toml"
	"os"
	"strings"
)

var (
	DB = New()
)

/**
*设置数据库连接
*@param diver string
 */
func New() *gorm.DB {

	if getAppEnv() == "test" {
		configTree := config.Conf.Get("test").(*toml.Tree)
		DB, err := gorm.Open(configTree.Get("DataBaseDriver").(string), configTree.Get("DataBaseConnect").(string))
		if err != nil {
			panic(fmt.Sprintf("No error should happen when connecting to  database, but got err=%+v", err))
		}

		return DB

	} else {

		driver := config.Conf.Get("database.dirver").(string)
		configTree := config.Conf.Get(driver).(*toml.Tree)
		userName := configTree.Get("databaseUserName").(string)
		password := configTree.Get("databasePassword").(string)
		databaseName := configTree.Get("databaseName").(string)
		connect := userName + ":" + password + "@/" + databaseName + "?charset=utf8&parseTime=True&loc=Local"

		fmt.Println(connect)

		DB, err := gorm.Open(driver, connect)

		if err != nil {
			panic(fmt.Sprintf("No error should happen when connecting to  database, but got err=%+v", err))
		}

		return DB
	}
}

//获取程序运行环境
// 根据程序运行路径后缀判断
//如果是 test 就是测试环境
func getAppEnv() string {
	file := os.Args[0]
	s := strings.Split(file, ".")
	return s[len(s)-1]
}

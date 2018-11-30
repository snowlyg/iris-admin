package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pelletier/go-toml"
)

/**
*设置数据库连接
*@param diver string
 */
func New(conf *toml.Tree, appEnv string) *gorm.DB {

	if appEnv == "testing" {

		configTree := conf.Get("test").(*toml.Tree)
		DB, err := gorm.Open(configTree.Get("DataBaseDriver").(string), configTree.Get("DataBaseConnect").(string))
		if err != nil {
			panic(fmt.Sprintf("No error should happen when connecting to  database, but got err=%+v", err))
		}

		return DB

	} else {

		driver := conf.Get("database.dirver").(string)
		configTree := conf.Get(driver).(*toml.Tree)
		DB, err := gorm.Open(driver, configTree.Get("connect").(string))

		if err != nil {
			panic(fmt.Sprintf("No error should happen when connecting to  database, but got err=%+v", err))
		}

		return DB
	}
}

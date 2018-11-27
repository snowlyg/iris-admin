package system

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
func setDatabase(diver string) {
	val := RedisGet("env_t")
	if val == "testing" {
		configTree := Config.Get("test").(*toml.Tree)
		DB, err = gorm.Open(configTree.Get("DataBaseDriver").(string), configTree.Get("DataBaseConnect").(string))
	} else {
		configTree := Config.Get(diver).(*toml.Tree)
		DB, err = gorm.Open(configTree.Get(diver).(string), configTree.Get("connect").(string))
	}

	if err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to  database, but got err=%+v", err))
	}
}

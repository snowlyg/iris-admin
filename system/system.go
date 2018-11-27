package system

import (
	"IrisYouQiKangApi/tools"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/pelletier/go-toml"
)

var (
	DB     *gorm.DB      //mysql
	Tools  *tools.Tools  //tools
	Config *toml.Tree    //config
	Redis  *redis.Client //redis
	err    error
)

func init() {
	Tools = tools.New()
	Config = configNew()

	//初始化reids
	newRedis()

	//设置数据库连接
	setDatabase(Config.Get("database.dirver").(string))
}

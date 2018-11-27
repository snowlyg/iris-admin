package system

import (
	"IrisYouQiKangApi/config"
	"IrisYouQiKangApi/tools"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pelletier/go-toml"
	"time"
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
	Config = config.New()
	configTree := Config.Get("redis").(*toml.Tree)
	Redis = redis.NewClient(&redis.Options{
		Addr:     configTree.Get("Addr").(string),
		Password: configTree.Get("Password").(string), // no password set
		DB:       int(configTree.Get("DB").(int64)),   // 因为系统是64位的，所以默认的 int 型是 int64
	})

	val := RedisGet("env_t")

	if val == "testing" {
		DB, err = gorm.Open("sqlite3", Config.Get("sqlite.connect").(string))
	} else {
		DB, err = gorm.Open("mysql", Config.Get("mysql.connect").(string))
	}

	if err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to  database, but got err=%+v", err))
	}
}

func RedisGet(key string) (value string) {
	value, err := Redis.Get(key).Result()
	if err == redis.Nil {
		fmt.Println("env_t does not exist")
	} else if err != nil {
		panic(err)
	}

	return
}

func RedisSet(key, value string, expiration time.Duration) {
	err = Redis.Set(key, value, expiration).Err()
	if err != nil {
		panic(err)
	}
}

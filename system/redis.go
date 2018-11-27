package system

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pelletier/go-toml"
	"time"
)

/**
*设置数据库连接
*@param diver string
 */
func newRedis() {
	configTree := Config.Get("redis").(*toml.Tree)
	Redis = redis.NewClient(&redis.Options{
		Addr:     configTree.Get("Addr").(string),
		Password: configTree.Get("Password").(string), // no password set
		DB:       int(configTree.Get("DB").(int64)),   // 因为系统是64位的，所以默认的 int 型是 int64
	})
}

//获取 reids 数据
func RedisGet(key string) (value string) {
	value, err := Redis.Get(key).Result()
	if err == redis.Nil {
		fmt.Println("env_t does not exist")
	} else if err != nil {
		panic(err)
	}

	return
}

//设置 reids 数据
func RedisSet(key, value string, expiration time.Duration) {
	err = Redis.Set(key, value, expiration).Err()
	if err != nil {
		panic(err)
	}
}

package logic

import (
	"IrisYouQiKangApi/config"
	"IrisYouQiKangApi/redis"
	"IrisYouQiKangApi/tools"
)

var (
	Tools  *tools.Tools
	Redis  *redis.Redis
	Config *config.Config //config
)

func init() {
	Config = config.New()
	Tools = tools.New()
	Redis = redis.New(Config.Redis.Connect, Config.Redis.DB, Config.Redis.MaxIdle, Config.Redis.MaxActive)
}

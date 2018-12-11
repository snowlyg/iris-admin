package caches

import (
	"IrisApiProject/config"
	"github.com/go-redis/redis"
	"github.com/pelletier/go-toml"
)

var (
	Cache = New()
)

func New() *redis.Client {
	configTree := config.Conf.Get("redis").(*toml.Tree)
	return redis.NewClient(&redis.Options{
		Addr:     configTree.Get("Addr").(string),
		Password: configTree.Get("Password").(string), // no password set
		DB:       int(configTree.Get("DB").(int64)),   // 因为系统是64位的，所以默认的 int 型是 int64
	})
}

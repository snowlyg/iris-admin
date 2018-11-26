package logic

import (
	"IrisYouQiKangApi/config"
	"github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
	"github.com/pelletier/go-toml"
	"time"
)

var (
	//Tools  *tools.Tools
	Redis  *redis.Database //redis
	Config *toml.Tree      //config
)

func init() {
	//Tools = tools.New()
	Config = config.New()
	configTree := Config.Get("redis").(*toml.Tree)
	Redis = redis.New(service.Config{
		Network:     configTree.Get("Network").(string),
		Addr:        configTree.Get("Addr").(string),
		Password:    configTree.Get("Password").(string),
		Database:    configTree.Get("Database").(string),
		MaxIdle:     int(configTree.Get("MaxIdle").(int64)),
		MaxActive:   int(configTree.Get("MaxActive").(int64)),
		IdleTimeout: time.Duration(5) * time.Minute,
		Prefix:      configTree.Get("Network").(string)}) // optionally configure the bridge between your redis server
}

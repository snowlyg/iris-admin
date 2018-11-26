package system

import (
	"IrisYouQiKangApi/config"
	"IrisYouQiKangApi/tools"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
	"github.com/pelletier/go-toml"
	"time"
)

var (
	DB     *gorm.DB        //mysql
	Tools  *tools.Tools    //tools
	Config *toml.Tree      //config
	Redis  *redis.Database //redis
	err    error
)

func init() {
	Tools = tools.New()
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
		Prefix:      configTree.Get("Network").(string)},
	)

	env := Redis.Get("iris", "my_env")

	Tools.Debug(env)

	DB, err = gorm.Open("mysql", Config.Get("mysql.connect").(string))
	if err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to  database, but got err=%+v", err))
	}
	//DB.AutoMigrate(&OauthToken{}, &Users{})
}

package models

import (
	"GoYouQiKangApi/models"
	"IrisYouQiKangApi/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	DB *xorm.Engine //mysql
	//Tools  *tools.Tools     //tools
	Config *config.Config //config
	//Modb   *mongodb.Mongodb //mongodb
	//Redis  *redis.Redis     //redis
)

func init() {
	Config = config.New()
	//Tools = tools.New()
	//Modb = mongodb.New(Config.Mongodb.Connect)
	//Redis = redis.New(Config.Redis.Connect, Config.Redis.DB, Config.Redis.MaxIdle, Config.Redis.MaxActive)

	var err error
	DB, err = xorm.NewEngine("mysql", Config.Mysql.Connect)
	if err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to  database, but got err=%+v", err))
	}
	err = DB.Sync2(new(models.Users), new(models.Roles))
	if err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to  database, but got err=%+v", err))
	}
}

package main

import (
	"IrisYouQiKangApi/config"
	"IrisYouQiKangApi/database"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/pelletier/go-toml"
)

var (
	db   *gorm.DB
	conf *toml.Tree
)

func main() {
	//初始化配置
	conf = config.New()
	//初始化测试数据库
	db = database.New(conf, "dev")

	app := NewApp()

	app.Run(iris.Addr(":80"))
}

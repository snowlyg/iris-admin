package main

import (
	_ "IrisYouQiKangApi/models"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris"
)

func init() {

	maxIdle := 30
	maxConn := 30
	dataSource := "root:UHC0JC5s6DEg9BRXYuDJnqbdl1ecL4gV@tcp(127.0.0.1:3306)/goyouqikang?charset=utf8&loc=Asia%2FShanghai"

	orm.RegisterDataBase("default", "mysql", dataSource, maxIdle, maxConn)

	// create table
	orm.RunSyncdb("default", false, true)

}

func main() {


	app := newApp()

	app.Run(iris.Addr(":8080"))

}

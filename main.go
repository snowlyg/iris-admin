package main

import (
	"IrisYouQiKangApi/routers"
	"github.com/kataras/iris"
)

func main() {
	myapp := new(routers.MyApp)
	app := myapp.NewApp()

	app.Run(iris.Addr(":80"))
}

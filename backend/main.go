package main

import (
	"fmt"
	"github.com/snowlyg/IrisAdminApi/backend/config"
	"time"

	"github.com/kataras/iris/v12"
)

func main() {

	f := NewLogFile()
	defer f.Close()

	app := NewApp()
	app.Logger().SetOutput(f) //记录日志

	if config.Config.HTTPS {
		host := fmt.Sprintf("%s:%d", config.Config.Host, 443)
		if err := app.Run(iris.TLS(host, config.Config.Certpath, config.Config.Certkey)); err != nil {
			fmt.Println(err)
		}
	} else {
		println(fmt.Sprintf("runing"))
		if err := app.Run(
			iris.Addr(fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port)),
			iris.WithoutServerError(iris.ErrServerClosed),
			iris.WithOptimizations,
			iris.WithTimeFormat(time.RFC3339),
		); err != nil {
			fmt.Println(err)
		}
	}
}

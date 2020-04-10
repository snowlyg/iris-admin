package main

import (
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/backend/config"
)

func main() {
	f := NewLogFile()
	defer f.Close()

	api := NewApp()
	//api.Logger().SetOutput(f) //记录日志

	if config.Config.HTTPS {
		host := fmt.Sprintf("%s:%d", config.Config.Host, 443)
		if err := api.Run(iris.TLS(host, config.Config.Certpath, config.Config.Certkey)); err != nil {
			fmt.Println(err)
		}
	} else {
		if err := api.Run(
			iris.Addr(fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port)),
			iris.WithoutServerError(iris.ErrServerClosed),
			iris.WithOptimizations,
			iris.WithTimeFormat(time.RFC3339),
		); err != nil {
			fmt.Println(err)
		}
	}

}

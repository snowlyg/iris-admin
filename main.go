package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/config"
)

func main() {
	f := NewLogFile()
	defer f.Close()

	api := NewApp()
	api.Logger().SetOutput(f) //记录日志

	if err := api.Run(iris.Addr(config.GetAppUrl()), iris.WithConfiguration(config.Isc)); err != nil {
		color.Yellow(fmt.Sprintf("项目运行结束: %v", err))
	}
}

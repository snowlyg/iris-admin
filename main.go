package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	irisbaserabc "github.com/snowlyg/iris-base-rabc"
	"github.com/snowlyg/iris-base-rabc/config"
)

func main() {
	f := irisbaserabc.NewLogFile()
	defer f.Close()

	api := irisbaserabc.NewApp()
	api.Logger().SetOutput(f) //记录日志
	err := api.Run(iris.Addr(config.GetAppUrl()), iris.WithConfiguration(config.GetIrisConf()))
	if err != nil {
		color.Yellow(fmt.Sprintf("项目运行结束: %v", err))
	}
}

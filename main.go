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

	// 修改系统默认设置方法
	//config.SetAppDriverType("Mysql")
	//config.SetMysqlConnect("root:password@(127.0.0.1:3306)/")
	//config.SetMysqlName("root")
	//config.SetAppUrl("http://localhost:8080")
	//config.SetAppName("MySiteName")
	//config.SetAppCreateSysData(false)
	//config.SetAppLoggerLevel("debug")
	//config.SetTestDataName("username")

	if err := api.Run(iris.Addr(config.GetAppUrl()), iris.WithConfiguration(config.GetIrisConf())); err != nil {
		color.Yellow(fmt.Sprintf("项目运行结束: %v", err))
	}
}

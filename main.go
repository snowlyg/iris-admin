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

	// 修改默认设置
	//config.SetAppDriverType("Mysql")
	//config.SetMysqlConnect("root:password@(127.0.0.1:3306)/")
	//config.SetMysqlName("root")
	//config.SetAppUrl("localhost:8081")
	//config.SetAppName("MySiteName")
	//config.SetAppCreateSysData(false)
	//config.SetAppLoggerLevel("debug")
	//config.SetTestDataName("username")

	url := config.GetAppUrl()
	conf := config.GetIrisConf()
	if err := api.Run(iris.Addr(url), iris.WithConfiguration(conf)); err != nil {
		color.Yellow(fmt.Sprintf("项目运行结束: %v", err))
	}
}
